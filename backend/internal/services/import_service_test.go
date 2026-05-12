package services

import (
	"errors"
	"testing"

	"rygell-dashboard/internal/models"

	"gorm.io/gorm"
)

type stubImportParser struct {
	sheets []ParsedSheet
	err    error
}

func (p stubImportParser) ParseExcelFile(string) ([]ParsedSheet, error) {
	return p.sheets, p.err
}

type fakeImportTxRunner struct {
	master     importMasterRepository
	contract   importContractRepository
	committed  bool
	rolledBack bool
}

func (r *fakeImportTxRunner) RunInTransaction(fn func(importMasterRepository, importContractRepository) error) error {
	err := fn(r.master, r.contract)
	if err != nil {
		r.rolledBack = true
		return err
	}
	r.committed = true
	return nil
}

type fakeImportMasterRepo struct {
	nextID uint
}

func (r *fakeImportMasterRepo) next() uint {
	r.nextID++
	return r.nextID
}

func (r *fakeImportMasterRepo) FindOrCreateVendorByName(name string) (*models.Vendor, error) {
	return &models.Vendor{ID: r.next(), Name: name, Code: name}, nil
}

func (r *fakeImportMasterRepo) FindOrCreateMillByName(name string) (*models.Mill, error) {
	return &models.Mill{ID: r.next(), Name: name, Code: name}, nil
}

func (r *fakeImportMasterRepo) FindOrCreateProductByName(name string) (*models.Product, error) {
	return &models.Product{ID: r.next(), Name: name}, nil
}

func (r *fakeImportMasterRepo) FindOrCreateZoneByName(name string) (*models.Zone, error) {
	return &models.Zone{ID: r.next(), Name: name}, nil
}

func (r *fakeImportMasterRepo) FindOrCreateMotByName(name string) (*models.Mot, error) {
	return &models.Mot{ID: r.next(), Name: name}, nil
}

func (r *fakeImportMasterRepo) FindOrCreateUomByName(name string) (*models.Uom, error) {
	return &models.Uom{ID: r.next(), Name: name}, nil
}

type fakeImportContractRepo struct {
	existingFix    map[string]bool
	existingVar    map[string]bool
	existingOncall map[string]bool
	createdFix     []models.ContractDedicatedFix
	createdVar     []models.ContractDedicatedVar
	createdOncall  []models.ContractOncall
}

func (r *fakeImportContractRepo) BulkCreateDedicatedFix(contracts []models.ContractDedicatedFix) error {
	r.createdFix = append(r.createdFix, contracts...)
	return nil
}

func (r *fakeImportContractRepo) BulkCreateDedicatedVar(contracts []models.ContractDedicatedVar) error {
	r.createdVar = append(r.createdVar, contracts...)
	return nil
}

func (r *fakeImportContractRepo) BulkCreateOncall(contracts []models.ContractOncall) error {
	r.createdOncall = append(r.createdOncall, contracts...)
	return nil
}

func (r *fakeImportContractRepo) FindDedicatedFixBySPKVendorMill(spk string, vendorID, millID uint) (*models.ContractDedicatedFix, error) {
	key, _ := contractDedupeKey(spk, vendorID, millID)
	if r.existingFix[key] {
		return &models.ContractDedicatedFix{ID: 1, SPKNumber: spk, VendorID: vendorID, MillID: millID}, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (r *fakeImportContractRepo) FindDedicatedVarBySPKVendorMill(spk string, vendorID, millID uint) (*models.ContractDedicatedVar, error) {
	key, _ := contractDedupeKey(spk, vendorID, millID)
	if r.existingVar[key] {
		return &models.ContractDedicatedVar{ID: 1, SPKNumber: spk, VendorID: vendorID, MillID: millID}, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (r *fakeImportContractRepo) FindOncallBySPKVendorMill(spk string, vendorID, millID uint) (*models.ContractOncall, error) {
	key, _ := contractDedupeKey(spk, vendorID, millID)
	if r.existingOncall[key] {
		return &models.ContractOncall{ID: 1, SPKNumber: spk, VendorID: vendorID, MillID: millID}, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func TestConfirmImportInsertsValidRows(t *testing.T) {
	master := &fakeImportMasterRepo{}
	contract := &fakeImportContractRepo{}
	tx := &fakeImportTxRunner{master: master, contract: contract}
	service := &ImportService{
		parserService: stubImportParser{sheets: []ParsedSheet{{
			SheetName: "Dedicated Var",
			SheetType: "dedicated_var",
			Rows: []map[string]string{{
				"TRANSPORTER/CARRIER": "Vendor A",
				"MILL/CATEGORY":       "Mill A",
				"SPK NUMBER":          "SPK-001",
				"PAYLOAD":             "10",
				"COST/KG":             "100",
			}},
		}}},
		masterRepo:   master,
		contractRepo: contract,
		txRunner:     tx,
	}

	result, err := service.ConfirmImport("ignored.xlsx")
	if err != nil {
		t.Fatalf("ConfirmImport() error = %v", err)
	}
	if !tx.committed || tx.rolledBack {
		t.Fatalf("transaction committed=%v rolledBack=%v, want committed only", tx.committed, tx.rolledBack)
	}
	if result.DedicatedVarInserted != 1 || len(contract.createdVar) != 1 {
		t.Fatalf("inserted=%d created=%d, want 1", result.DedicatedVarInserted, len(contract.createdVar))
	}
}

func TestConfirmParsedSheetsImportsEditedPreviewRows(t *testing.T) {
	master := &fakeImportMasterRepo{}
	contract := &fakeImportContractRepo{}
	tx := &fakeImportTxRunner{master: master, contract: contract}
	service := &ImportService{
		parserService: stubImportParser{err: errors.New("parser should not be called")},
		masterRepo:    master,
		contractRepo:  contract,
		txRunner:      tx,
	}

	result, err := service.ConfirmParsedSheets([]ParsedSheet{{
		SheetName: "Edited Dedicated Var",
		SheetType: "dedicated_var",
		Rows: []map[string]string{{
			"TRANSPORTER/CARRIER": "Edited Vendor",
			"MILL/CATEGORY":       "Edited Mill",
			"SPK NUMBER":          "SPK-EDITED",
			"PAYLOAD":             "12",
			"COST/KG":             "150",
		}},
	}})
	if err != nil {
		t.Fatalf("ConfirmParsedSheets() error = %v", err)
	}
	if !tx.committed || tx.rolledBack {
		t.Fatalf("transaction committed=%v rolledBack=%v, want committed only", tx.committed, tx.rolledBack)
	}
	if result.DedicatedVarInserted != 1 || len(contract.createdVar) != 1 {
		t.Fatalf("inserted=%d created=%d, want 1", result.DedicatedVarInserted, len(contract.createdVar))
	}
	if contract.createdVar[0].SPKNumber != "SPK-EDITED" {
		t.Fatalf("SPKNumber = %q, want edited value", contract.createdVar[0].SPKNumber)
	}
}

func TestConfirmImportSkipsDuplicateRows(t *testing.T) {
	master := &fakeImportMasterRepo{}
	contract := &fakeImportContractRepo{existingFix: map[string]bool{
		"spk-001|vendor:1|mill:2": true,
	}}
	tx := &fakeImportTxRunner{master: master, contract: contract}
	service := &ImportService{
		parserService: stubImportParser{sheets: []ParsedSheet{{
			SheetName: "Dedicated Fix",
			SheetType: "dedicated_fix",
			Rows: []map[string]string{{
				"TRANSPORTER/CARRIER": "Vendor A",
				"MILL/CATEGORY":       "Mill A",
				"SPK NUMBER":          "SPK-001",
			}},
		}}},
		masterRepo:   master,
		contractRepo: contract,
		txRunner:     tx,
	}

	result, err := service.ConfirmImport("ignored.xlsx")
	if err != nil {
		t.Fatalf("ConfirmImport() error = %v", err)
	}
	if result.DedicatedFixInserted != 0 || result.DedicatedFixSkippedDuplicates != 1 {
		t.Fatalf("inserted=%d skipped=%d, want inserted=0 skipped=1", result.DedicatedFixInserted, result.DedicatedFixSkippedDuplicates)
	}
	if len(contract.createdFix) != 0 {
		t.Fatalf("created %d duplicate rows, want 0", len(contract.createdFix))
	}
}

func TestConfirmImportRollsBackOnRowErrors(t *testing.T) {
	master := &fakeImportMasterRepo{}
	contract := &fakeImportContractRepo{}
	tx := &fakeImportTxRunner{master: master, contract: contract}
	service := &ImportService{
		parserService: stubImportParser{sheets: []ParsedSheet{{
			SheetName: "Oncall",
			SheetType: "oncall",
			Rows: []map[string]string{{
				"TRANSPORTER/CARRIER": "Vendor A",
				"MILL/CATEGORY":       "Mill A",
			}},
		}}},
		masterRepo:   master,
		contractRepo: contract,
		txRunner:     tx,
	}

	result, err := service.ConfirmImport("ignored.xlsx")
	if err == nil {
		t.Fatal("ConfirmImport() error = nil, want validation error")
	}
	var validationErr *ImportValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("ConfirmImport() error = %T, want ImportValidationError", err)
	}
	if !tx.rolledBack || tx.committed {
		t.Fatalf("transaction committed=%v rolledBack=%v, want rollback only", tx.committed, tx.rolledBack)
	}
	if len(result.Errors) == 0 {
		t.Fatal("result.Errors is empty, want row error")
	}
}
