package services

import (
	"encoding/json"
	"fmt"
	"strconv"

	"rygell-dashboard/internal/models"
	"rygell-dashboard/internal/repositories"

	"gorm.io/gorm"
)

// ContractService handles business logic for contracts and audit logging.
type ContractService struct {
	contractRepo *repositories.ContractRepository
	auditRepo    auditLogStore
	masterRepo   agreementMasterStore
	txRunner     contractTransactionRunner
}

type auditLogStore interface {
	Create(log *models.AuditLog) error
	GetByEntity(entityType string, entityID uint) ([]models.AuditLog, error)
	GetByVendor(vendorID uint) ([]models.AuditLog, error)
}

type agreementMasterStore interface {
	GetVendorByID(id uint) (*models.Vendor, error)
	GetMillByID(id uint) (*models.Mill, error)
}

type contractTransactionRunner interface {
	RunInTransaction(func(*repositories.ContractRepository, auditLogStore, agreementMasterStore) error) error
}

type gormContractTransactionRunner struct {
	db         *gorm.DB
	withMaster bool
}

func (r *gormContractTransactionRunner) RunInTransaction(fn func(*repositories.ContractRepository, auditLogStore, agreementMasterStore) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var masterRepo agreementMasterStore
		if r.withMaster {
			masterRepo = repositories.NewMasterRepository(tx)
		}
		return fn(
			repositories.NewContractRepository(tx),
			repositories.NewAuditRepository(tx),
			masterRepo,
		)
	})
}

// NewContractService creates a new ContractService.
func NewContractService(contractRepo *repositories.ContractRepository, auditRepo *repositories.AuditRepository, masterRepo ...*repositories.MasterRepository) *ContractService {
	service := &ContractService{
		contractRepo: contractRepo,
		auditRepo:    auditRepo,
	}
	if len(masterRepo) > 0 {
		service.masterRepo = masterRepo[0]
	}
	if contractRepo != nil && auditRepo != nil {
		service.txRunner = &gormContractTransactionRunner{
			db:         contractRepo.DB(),
			withMaster: len(masterRepo) > 0,
		}
	}
	return service
}

// --- Dedicated Fix ---

func (s *ContractService) GetAllDedicatedFix(filters map[string]interface{}, search string) ([]models.ContractDedicatedFix, error) {
	return s.contractRepo.GetAllDedicatedFix(filters, search)
}

func (s *ContractService) GetDedicatedFixPage(filters map[string]interface{}, search string, limit, offset int) ([]models.ContractDedicatedFix, int64, error) {
	return s.contractRepo.GetDedicatedFixPage(filters, search, limit, offset)
}

func (s *ContractService) GetDedicatedFixByID(id uint) (*models.ContractDedicatedFix, error) {
	return s.contractRepo.GetDedicatedFixByID(id)
}

func (s *ContractService) CreateDedicatedFix(contract *models.ContractDedicatedFix) error {
	return s.withMutationTransaction(func(worker *ContractService) error {
		if err := worker.contractRepo.CreateDedicatedFix(contract); err != nil {
			return err
		}
		return worker.logAudit("contract_dedicated_fix", contract.ID, "create", "", "", nil, contract)
	})
}

func (s *ContractService) UpdateDedicatedFix(contract *models.ContractDedicatedFix, changedBy, note string) error {
	return s.withMutationTransaction(func(worker *ContractService) error {
		old, err := worker.contractRepo.GetDedicatedFixByID(contract.ID)
		if err != nil {
			return err
		}
		if err := worker.contractRepo.UpdateDedicatedFix(contract); err != nil {
			return err
		}
		updated, err := worker.contractRepo.GetDedicatedFixByID(contract.ID)
		if err != nil {
			return err
		}
		return worker.logAudit("contract_dedicated_fix", contract.ID, "update", changedBy, note, old, updated)
	})
}

func (s *ContractService) DeleteDedicatedFix(id uint) error {
	return s.contractRepo.DeleteDedicatedFix(id)
}

// UpdateDedicatedFixAgreement updates only the agreement note on a contract.
func (s *ContractService) UpdateDedicatedFixAgreement(id uint, changedBy, note string) error {
	return s.withMutationTransaction(func(worker *ContractService) error {
		contract, err := worker.contractRepo.GetDedicatedFixByID(id)
		if err != nil {
			return err
		}
		old := *contract
		contract.Notes = note
		if err := worker.contractRepo.UpdateDedicatedFix(contract); err != nil {
			return err
		}
		updated, err := worker.contractRepo.GetDedicatedFixByID(id)
		if err != nil {
			return err
		}
		return worker.logAudit("contract_dedicated_fix", id, "agreement_update", changedBy, note, &old, updated)
	})
}

// --- Dedicated Var ---

func (s *ContractService) GetAllDedicatedVar(filters map[string]interface{}, search string) ([]models.ContractDedicatedVar, error) {
	return s.contractRepo.GetAllDedicatedVar(filters, search)
}

func (s *ContractService) GetDedicatedVarPage(filters map[string]interface{}, search string, limit, offset int) ([]models.ContractDedicatedVar, int64, error) {
	return s.contractRepo.GetDedicatedVarPage(filters, search, limit, offset)
}

func (s *ContractService) GetDedicatedVarByID(id uint) (*models.ContractDedicatedVar, error) {
	return s.contractRepo.GetDedicatedVarByID(id)
}

func (s *ContractService) CreateDedicatedVar(contract *models.ContractDedicatedVar) error {
	return s.withMutationTransaction(func(worker *ContractService) error {
		if err := worker.contractRepo.CreateDedicatedVar(contract); err != nil {
			return err
		}
		return worker.logAudit("contract_dedicated_var", contract.ID, "create", "", "", nil, contract)
	})
}

func (s *ContractService) UpdateDedicatedVar(contract *models.ContractDedicatedVar, changedBy, note string) error {
	return s.withMutationTransaction(func(worker *ContractService) error {
		old, err := worker.contractRepo.GetDedicatedVarByID(contract.ID)
		if err != nil {
			return err
		}
		if err := worker.contractRepo.UpdateDedicatedVar(contract); err != nil {
			return err
		}
		updated, err := worker.contractRepo.GetDedicatedVarByID(contract.ID)
		if err != nil {
			return err
		}
		return worker.logAudit("contract_dedicated_var", contract.ID, "update", changedBy, note, old, updated)
	})
}

func (s *ContractService) DeleteDedicatedVar(id uint) error {
	return s.contractRepo.DeleteDedicatedVar(id)
}

// UpdateDedicatedVarAgreement updates only the agreement note on a var contract.
func (s *ContractService) UpdateDedicatedVarAgreement(id uint, changedBy, note string) error {
	return s.withMutationTransaction(func(worker *ContractService) error {
		contract, err := worker.contractRepo.GetDedicatedVarByID(id)
		if err != nil {
			return err
		}
		old := *contract
		contract.Notes = note
		if err := worker.contractRepo.UpdateDedicatedVar(contract); err != nil {
			return err
		}
		updated, err := worker.contractRepo.GetDedicatedVarByID(id)
		if err != nil {
			return err
		}
		return worker.logAudit("contract_dedicated_var", id, "agreement_update", changedBy, note, &old, updated)
	})
}

// --- Oncall ---

func (s *ContractService) GetAllOncall(filters map[string]interface{}, search string) ([]models.ContractOncall, error) {
	return s.contractRepo.GetAllOncall(filters, search)
}

func (s *ContractService) GetOncallPage(filters map[string]interface{}, search string, limit, offset int) ([]models.ContractOncall, int64, error) {
	return s.contractRepo.GetOncallPage(filters, search, limit, offset)
}

func (s *ContractService) GetOncallByID(id uint) (*models.ContractOncall, error) {
	return s.contractRepo.GetOncallByID(id)
}

func (s *ContractService) CreateOncall(contract *models.ContractOncall) error {
	return s.withMutationTransaction(func(worker *ContractService) error {
		if err := worker.contractRepo.CreateOncall(contract); err != nil {
			return err
		}
		return worker.logAudit("contract_oncall", contract.ID, "create", "", "", nil, contract)
	})
}

func (s *ContractService) UpdateOncall(contract *models.ContractOncall, changedBy, note string) error {
	return s.withMutationTransaction(func(worker *ContractService) error {
		old, err := worker.contractRepo.GetOncallByID(contract.ID)
		if err != nil {
			return err
		}
		if err := worker.contractRepo.UpdateOncall(contract); err != nil {
			return err
		}
		updated, err := worker.contractRepo.GetOncallByID(contract.ID)
		if err != nil {
			return err
		}
		return worker.logAudit("contract_oncall", contract.ID, "update", changedBy, note, old, updated)
	})
}

func (s *ContractService) DeleteOncall(id uint) error {
	return s.contractRepo.DeleteOncall(id)
}

// UpdateOncallAgreement updates only the agreement note on an oncall contract.
func (s *ContractService) UpdateOncallAgreement(id uint, changedBy, note string) error {
	return s.withMutationTransaction(func(worker *ContractService) error {
		contract, err := worker.contractRepo.GetOncallByID(id)
		if err != nil {
			return err
		}
		old := *contract
		contract.Notes = note
		if err := worker.contractRepo.UpdateOncall(contract); err != nil {
			return err
		}
		updated, err := worker.contractRepo.GetOncallByID(id)
		if err != nil {
			return err
		}
		return worker.logAudit("contract_oncall", id, "agreement_update", changedBy, note, &old, updated)
	})
}

// --- Map-based Updates (partial update, no association conflicts) ---

// sanitizeUpdateMap removes nested objects and metadata fields from the update map
// so GORM only updates scalar columns.
func sanitizeUpdateMap(body map[string]interface{}) map[string]interface{} {
	// Keys to delete: nested objects and metadata
	deleteKeys := []string{
		"mill", "vendor", "product", "mot", "uom",
		"origin_zone", "dest_zone", "created_at", "updated_at", "deleted_at", "id",
	}
	for _, k := range deleteKeys {
		delete(body, k)
	}

	// Convert string numeric values to float64 for decimal columns
	numericKeys := []string{
		"distance", "payload", "loading_cost", "unloading_cost",
		"cost_idr", "cost_per_kg", "cost_per_ton", "cost_per_kg_km",
		"running_cost_idr", "running_cost_usd",
		"cost_jan", "cost_feb", "cost_mar", "cost_apr", "cost_may", "cost_jun",
		"fix_cost", "distributed_cost", "cargo_carried", "unit_cost",
	}
	for _, k := range numericKeys {
		if v, ok := body[k]; ok {
			switch val := v.(type) {
			case string:
				if val == "" {
					body[k] = 0
				} else {
					parsed, err := strconv.ParseFloat(val, 64)
					if err == nil {
						body[k] = parsed
					} else {
						body[k] = 0
					}
				}
			}
		}
	}
	return body
}

func (s *ContractService) UpdateDedicatedFixMap(id uint, body map[string]interface{}, changedBy, note string) (*models.ContractDedicatedFix, error) {
	var result *models.ContractDedicatedFix
	err := s.withMutationTransaction(func(worker *ContractService) error {
		old, err := worker.contractRepo.GetDedicatedFixByID(id)
		if err != nil {
			return err
		}
		updates := sanitizeUpdateMap(body)
		if err := worker.contractRepo.UpdateDedicatedFixMap(id, updates); err != nil {
			return err
		}
		result, err = worker.contractRepo.GetDedicatedFixByID(id)
		if err != nil {
			return err
		}
		return worker.logAudit("contract_dedicated_fix", id, "update", changedBy, note, old, result)
	})
	return result, err
}

func (s *ContractService) UpdateDedicatedVarMap(id uint, body map[string]interface{}, changedBy, note string) (*models.ContractDedicatedVar, error) {
	var result *models.ContractDedicatedVar
	err := s.withMutationTransaction(func(worker *ContractService) error {
		old, err := worker.contractRepo.GetDedicatedVarByID(id)
		if err != nil {
			return err
		}
		updates := sanitizeUpdateMap(body)
		if err := worker.contractRepo.UpdateDedicatedVarMap(id, updates); err != nil {
			return err
		}
		result, err = worker.contractRepo.GetDedicatedVarByID(id)
		if err != nil {
			return err
		}
		return worker.logAudit("contract_dedicated_var", id, "update", changedBy, note, old, result)
	})
	return result, err
}

func (s *ContractService) UpdateOncallMap(id uint, body map[string]interface{}, changedBy, note string) (*models.ContractOncall, error) {
	var result *models.ContractOncall
	err := s.withMutationTransaction(func(worker *ContractService) error {
		old, err := worker.contractRepo.GetOncallByID(id)
		if err != nil {
			return err
		}
		updates := sanitizeUpdateMap(body)
		if err := worker.contractRepo.UpdateOncallMap(id, updates); err != nil {
			return err
		}
		result, err = worker.contractRepo.GetOncallByID(id)
		if err != nil {
			return err
		}
		return worker.logAudit("contract_oncall", id, "update", changedBy, note, old, result)
	})
	return result, err
}

// --- Audit History ---

func (s *ContractService) GetAuditHistory(entityType string, entityID uint) ([]models.AuditLog, error) {
	return s.auditRepo.GetByEntity(entityType, entityID)
}

func (s *ContractService) GetVendorAuditHistory(vendorID uint) ([]models.AuditLog, error) {
	return s.auditRepo.GetByVendor(vendorID)
}

// UpdateVendorAgreement adds a negotiation note directly to a vendor.
func (s *ContractService) UpdateVendorAgreement(vendorID uint, changedBy, note string) error {
	return s.withMutationTransaction(func(worker *ContractService) error {
		if worker.masterRepo == nil {
			return fmt.Errorf("master repository is required for vendor agreement updates")
		}
		vendor, err := worker.masterRepo.GetVendorByID(vendorID)
		if err != nil {
			return err
		}
		return worker.logAudit("vendor", vendorID, "agreement_update", changedBy, note, vendor, vendor)
	})
}

// UpdateMillAgreement adds a negotiation note directly to a mill.
func (s *ContractService) UpdateMillAgreement(millID uint, changedBy, note string) error {
	return s.withMutationTransaction(func(worker *ContractService) error {
		if worker.masterRepo == nil {
			return fmt.Errorf("master repository is required for mill agreement updates")
		}
		mill, err := worker.masterRepo.GetMillByID(millID)
		if err != nil {
			return err
		}
		return worker.logAudit("mill", millID, "agreement_update", changedBy, note, mill, mill)
	})
}

// --- Internal Helpers ---

func (s *ContractService) withMutationTransaction(fn func(*ContractService) error) error {
	if s.txRunner == nil {
		return fn(s)
	}
	return s.txRunner.RunInTransaction(func(contractRepo *repositories.ContractRepository, auditRepo auditLogStore, masterRepo agreementMasterStore) error {
		worker := &ContractService{
			contractRepo: contractRepo,
			auditRepo:    auditRepo,
			masterRepo:   s.masterRepo,
		}
		if masterRepo != nil {
			worker.masterRepo = masterRepo
		}
		return fn(worker)
	})
}

func (s *ContractService) logAudit(entityType string, entityID uint, action, changedBy, note string, oldData, newData interface{}) error {
	audit := &models.AuditLog{
		EntityType:    entityType,
		EntityID:      entityID,
		Action:        action,
		ChangedBy:     changedBy,
		AgreementNote: note,
	}
	if oldData != nil {
		data, err := marshalAuditData(oldData)
		if err != nil {
			return err
		}
		audit.OldData = data
	}
	if newData != nil {
		data, err := marshalAuditData(newData)
		if err != nil {
			return err
		}
		audit.NewData = data
	}
	return s.auditRepo.Create(audit)
}

func marshalAuditData(value interface{}) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("failed to marshal audit data: %w", err)
	}
	return string(data), nil
}
