package services

import (
	"strings"
	"testing"

	"rygell-dashboard/internal/models"
	"rygell-dashboard/internal/repositories"

	"gorm.io/gorm"
)

type fakeAuditStore struct {
	logs []models.AuditLog
	err  error
}

type fakeAgreementMasterStore struct {
	vendors map[uint]*models.Vendor
	mills   map[uint]*models.Mill
}

type fakeContractTransactionRunner struct {
	called bool
	audits auditLogStore
	master agreementMasterStore
}

func (r *fakeContractTransactionRunner) RunInTransaction(fn func(*repositories.ContractRepository, auditLogStore, agreementMasterStore) error) error {
	r.called = true
	return fn(nil, r.audits, r.master)
}

func (s *fakeAgreementMasterStore) GetVendorByID(id uint) (*models.Vendor, error) {
	if vendor, ok := s.vendors[id]; ok {
		return vendor, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (s *fakeAgreementMasterStore) GetMillByID(id uint) (*models.Mill, error) {
	if mill, ok := s.mills[id]; ok {
		return mill, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (s *fakeAuditStore) Create(log *models.AuditLog) error {
	if s.err != nil {
		return s.err
	}
	s.logs = append(s.logs, *log)
	return nil
}

func (s *fakeAuditStore) GetByEntity(entityType string, entityID uint) ([]models.AuditLog, error) {
	if s.err != nil {
		return nil, s.err
	}
	var result []models.AuditLog
	for _, log := range s.logs {
		if log.EntityType == entityType && log.EntityID == entityID {
			result = append(result, log)
		}
	}
	return result, nil
}

func (s *fakeAuditStore) GetByVendor(vendorID uint) ([]models.AuditLog, error) {
	if s.err != nil {
		return nil, s.err
	}
	var result []models.AuditLog
	for _, log := range s.logs {
		if log.EntityType == "vendor" && log.EntityID == vendorID {
			result = append(result, log)
		}
	}
	return result, nil
}

func TestLogAuditStoresOldAndNewData(t *testing.T) {
	audits := &fakeAuditStore{}
	service := &ContractService{auditRepo: audits}

	oldData := map[string]string{"notes": "old"}
	newData := map[string]string{"notes": "new"}
	if err := service.logAudit("contract_oncall", 7, "update", "tester", "changed", oldData, newData); err != nil {
		t.Fatalf("logAudit() error = %v", err)
	}

	if len(audits.logs) != 1 {
		t.Fatalf("stored logs = %d, want 1", len(audits.logs))
	}
	log := audits.logs[0]
	if log.EntityType != "contract_oncall" || log.EntityID != 7 || log.Action != "update" {
		t.Fatalf("unexpected audit identity: %+v", log)
	}
	if !strings.Contains(log.OldData, `"notes":"old"`) {
		t.Fatalf("old_data = %q, want old notes payload", log.OldData)
	}
	if !strings.Contains(log.NewData, `"notes":"new"`) {
		t.Fatalf("new_data = %q, want new notes payload", log.NewData)
	}
}

func TestGetAuditHistoryRetrievesEntityLogs(t *testing.T) {
	audits := &fakeAuditStore{logs: []models.AuditLog{
		{EntityType: "contract_oncall", EntityID: 7, Action: "update"},
		{EntityType: "contract_oncall", EntityID: 8, Action: "update"},
		{EntityType: "vendor", EntityID: 7, Action: "agreement_update"},
	}}
	service := &ContractService{auditRepo: audits}

	logs, err := service.GetAuditHistory("contract_oncall", 7)
	if err != nil {
		t.Fatalf("GetAuditHistory() error = %v", err)
	}
	if len(logs) != 1 {
		t.Fatalf("logs = %d, want 1", len(logs))
	}
	if logs[0].EntityType != "contract_oncall" || logs[0].EntityID != 7 {
		t.Fatalf("unexpected log returned: %+v", logs[0])
	}
}

func TestMutationTransactionUsesTransactionScopedStores(t *testing.T) {
	transactionAudits := &fakeAuditStore{}
	runner := &fakeContractTransactionRunner{audits: transactionAudits}
	service := &ContractService{
		auditRepo: &fakeAuditStore{},
		txRunner:  runner,
	}

	err := service.withMutationTransaction(func(worker *ContractService) error {
		return worker.logAudit("contract_oncall", 9, "update", "tester", "tx", nil, map[string]string{"notes": "tx"})
	})
	if err != nil {
		t.Fatalf("withMutationTransaction() error = %v", err)
	}
	if !runner.called {
		t.Fatal("transaction runner was not called")
	}
	if len(transactionAudits.logs) != 1 {
		t.Fatalf("transaction audit logs = %d, want 1", len(transactionAudits.logs))
	}
}

func TestVendorAgreementRequiresExistingVendor(t *testing.T) {
	audits := &fakeAuditStore{}
	service := &ContractService{
		auditRepo: audits,
		masterRepo: &fakeAgreementMasterStore{
			vendors: map[uint]*models.Vendor{
				7: {ID: 7, Name: "Vendor A", Code: "V-A"},
			},
		},
	}

	if err := service.UpdateVendorAgreement(7, "tester", "confirmed"); err != nil {
		t.Fatalf("UpdateVendorAgreement() error = %v", err)
	}
	if len(audits.logs) != 1 {
		t.Fatalf("stored logs = %d, want 1", len(audits.logs))
	}
	if audits.logs[0].OldData == "" || audits.logs[0].NewData == "" {
		t.Fatalf("old/new data should be populated: %+v", audits.logs[0])
	}

	if err := service.UpdateVendorAgreement(999, "tester", "missing"); err == nil {
		t.Fatal("UpdateVendorAgreement() error = nil, want not found")
	}
	if len(audits.logs) != 1 {
		t.Fatalf("stored logs after missing vendor = %d, want 1", len(audits.logs))
	}
}
