package repositories

import (
	"strings"

	"rygell-dashboard/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ContractRepository handles CRUD for all contract types.
type ContractRepository struct {
	db *gorm.DB
}

// NewContractRepository creates a new ContractRepository.
func NewContractRepository(db *gorm.DB) *ContractRepository {
	return &ContractRepository{db: db}
}

// DB exposes the current GORM handle for transaction-scoped services.
func (r *ContractRepository) DB() *gorm.DB {
	return r.db
}

func updateExistingContract(db *gorm.DB, model interface{}, id uint, values interface{}) error {
	if err := db.First(model, id).Error; err != nil {
		return err
	}
	return db.Model(model).
		Select("*").
		Omit(clause.Associations, "id", "created_at", "deleted_at").
		Updates(values).Error
}

func updateExistingContractMap(db *gorm.DB, model interface{}, id uint, updates map[string]interface{}) error {
	if err := db.First(model, id).Error; err != nil {
		return err
	}
	if len(updates) == 0 {
		return nil
	}
	result := db.Model(model).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func deleteExistingContract(db *gorm.DB, model interface{}, id uint) error {
	result := db.Delete(model, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func normalizeSPK(spk string) string {
	return strings.ToLower(strings.TrimSpace(spk))
}

func applyCommonContractFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	if v, ok := filters["vendor_id"]; ok && v != "" {
		query = query.Where("vendor_id = ?", v)
	}
	if v, ok := filters["mill_id"]; ok && v != "" {
		query = query.Where("mill_id = ?", v)
	}
	return query
}

func (r *ContractRepository) buildDedicatedFixQuery(filters map[string]interface{}, search string) *gorm.DB {
	query := r.db.Model(&models.ContractDedicatedFix{}).Preload("Mill").Preload("Vendor").Preload("Product")
	query = applyCommonContractFilters(query, filters)
	if search == "" {
		return query
	}
	like := "%" + search + "%"
	return query.Select("DISTINCT contract_dedicated_fixes.*").
		Joins("LEFT JOIN vendors ON vendors.id = contract_dedicated_fixes.vendor_id").
		Joins("LEFT JOIN mills ON mills.id = contract_dedicated_fixes.mill_id").
		Joins("LEFT JOIN products ON products.id = contract_dedicated_fixes.product_id").
		Joins("LEFT JOIN mots ON mots.id = contract_dedicated_fixes.mot_id").
		Joins("LEFT JOIN uoms ON uoms.id = contract_dedicated_fixes.uom_id").
		Where(`
			contract_dedicated_fixes.spk_number ILIKE ? OR
			contract_dedicated_fixes.fa_number ILIKE ? OR
			contract_dedicated_fixes.area_category ILIKE ? OR
			contract_dedicated_fixes.proposal_cfas ILIKE ? OR
			contract_dedicated_fixes.license_plate ILIKE ? OR
			contract_dedicated_fixes.notes ILIKE ? OR
			vendors.name ILIKE ? OR
			vendors.code ILIKE ? OR
			mills.name ILIKE ? OR
			mills.code ILIKE ? OR
			products.name ILIKE ? OR
			mots.name ILIKE ? OR
			uoms.name ILIKE ? OR
			CAST(contract_dedicated_fixes.fix_cost AS TEXT) ILIKE ? OR
			CAST(contract_dedicated_fixes.distributed_cost AS TEXT) ILIKE ? OR
			CAST(contract_dedicated_fixes.unit_cost AS TEXT) ILIKE ? OR
			CAST(contract_dedicated_fixes.cost_per_kg AS TEXT) ILIKE ? OR
			CAST(contract_dedicated_fixes.cost_per_kgkm AS TEXT) ILIKE ? OR
			CAST(contract_dedicated_fixes.cargo_carried AS TEXT) ILIKE ?
		`,
			like, like, like, like, like, like,
			like, like, like, like, like, like, like,
			like, like, like, like, like, like,
		)
}

// --- Dedicated Fix ---

func (r *ContractRepository) GetAllDedicatedFix(filters map[string]interface{}, search string) ([]models.ContractDedicatedFix, error) {
	var contracts []models.ContractDedicatedFix
	query := r.buildDedicatedFixQuery(filters, search)
	err := query.Order("contract_dedicated_fixes.id ASC").Find(&contracts).Error
	return contracts, err
}

func (r *ContractRepository) GetDedicatedFixPage(filters map[string]interface{}, search string, limit, offset int) ([]models.ContractDedicatedFix, int64, error) {
	var contracts []models.ContractDedicatedFix
	query := r.buildDedicatedFixQuery(filters, search)
	countQuery := query.Session(&gorm.Session{})
	var total int64
	if err := countQuery.Distinct("contract_dedicated_fixes.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	if err := query.Order("contract_dedicated_fixes.id ASC").Find(&contracts).Error; err != nil {
		return nil, 0, err
	}
	return contracts, total, nil
}

func (r *ContractRepository) GetDedicatedFixByID(id uint) (*models.ContractDedicatedFix, error) {
	var contract models.ContractDedicatedFix
	err := r.db.Preload("Mill").Preload("Vendor").Preload("Product").First(&contract, id).Error
	return &contract, err
}

func (r *ContractRepository) CreateDedicatedFix(contract *models.ContractDedicatedFix) error {
	return r.db.Omit(clause.Associations).Create(contract).Error
}

func (r *ContractRepository) UpdateDedicatedFix(contract *models.ContractDedicatedFix) error {
	return updateExistingContract(r.db, &models.ContractDedicatedFix{}, contract.ID, contract)
}

func (r *ContractRepository) DeleteDedicatedFix(id uint) error {
	return deleteExistingContract(r.db, &models.ContractDedicatedFix{}, id)
}

func (r *ContractRepository) BulkCreateDedicatedFix(contracts []models.ContractDedicatedFix) error {
	return r.db.CreateInBatches(contracts, 100).Error
}

func (r *ContractRepository) FindDedicatedFixBySPK(spk string) (*models.ContractDedicatedFix, error) {
	var contract models.ContractDedicatedFix
	err := r.db.Where("spk_number = ?", spk).First(&contract).Error
	return &contract, err
}

func (r *ContractRepository) FindDedicatedFixBySPKVendorMill(spk string, vendorID, millID uint) (*models.ContractDedicatedFix, error) {
	normalized := normalizeSPK(spk)
	if normalized == "" {
		return nil, gorm.ErrRecordNotFound
	}
	var contract models.ContractDedicatedFix
	err := r.db.Where("LOWER(TRIM(spk_number)) = ? AND vendor_id = ? AND mill_id = ?", normalized, vendorID, millID).First(&contract).Error
	return &contract, err
}

// --- Dedicated Var ---

func (r *ContractRepository) GetAllDedicatedVar(filters map[string]interface{}, search string) ([]models.ContractDedicatedVar, error) {
	var contracts []models.ContractDedicatedVar
	query := r.db.Model(&models.ContractDedicatedVar{}).Preload("Mill").Preload("Vendor").Preload("Product").
		Preload("OriginZone").Preload("DestZone").Preload("Mot").Preload("Uom")
	query = applyCommonContractFilters(query, filters)
	if search != "" {
		like := "%" + search + "%"
		query = query.Select("DISTINCT contract_dedicated_vars.*").
			Joins("LEFT JOIN vendors ON vendors.id = contract_dedicated_vars.vendor_id").
			Joins("LEFT JOIN mills ON mills.id = contract_dedicated_vars.mill_id").
			Joins("LEFT JOIN products ON products.id = contract_dedicated_vars.product_id").
			Joins("LEFT JOIN zones AS origin_zones ON origin_zones.id = contract_dedicated_vars.origin_zone_id").
			Joins("LEFT JOIN zones AS dest_zones ON dest_zones.id = contract_dedicated_vars.dest_zone_id").
			Joins("LEFT JOIN mots ON mots.id = contract_dedicated_vars.mot_id").
			Joins("LEFT JOIN uoms ON uoms.id = contract_dedicated_vars.uom_id").
			Where(`
				contract_dedicated_vars.spk_number ILIKE ? OR
				contract_dedicated_vars.fa_number ILIKE ? OR
				contract_dedicated_vars.area_category ILIKE ? OR
				contract_dedicated_vars.proposal_cfas ILIKE ? OR
				contract_dedicated_vars.notes ILIKE ? OR
				vendors.name ILIKE ? OR
				vendors.code ILIKE ? OR
				mills.name ILIKE ? OR
				mills.code ILIKE ? OR
				products.name ILIKE ? OR
				origin_zones.name ILIKE ? OR
				dest_zones.name ILIKE ? OR
				mots.name ILIKE ? OR
				uoms.name ILIKE ? OR
				CAST(contract_dedicated_vars.distance AS TEXT) ILIKE ? OR
				CAST(contract_dedicated_vars.payload AS TEXT) ILIKE ? OR
				CAST(contract_dedicated_vars.cost_id_r AS TEXT) ILIKE ? OR
				CAST(contract_dedicated_vars.cost_per_kg AS TEXT) ILIKE ? OR
				CAST(contract_dedicated_vars.cost_per_kgkm AS TEXT) ILIKE ?
			`,
				like, like, like, like, like,
				like, like, like, like, like, like, like, like, like,
				like, like, like, like, like,
			)
	}
	err := query.Order("contract_dedicated_vars.id ASC").Find(&contracts).Error
	return contracts, err
}

func (r *ContractRepository) GetDedicatedVarPage(filters map[string]interface{}, search string, limit, offset int) ([]models.ContractDedicatedVar, int64, error) {
	var contracts []models.ContractDedicatedVar
	query := r.db.Model(&models.ContractDedicatedVar{}).Preload("Mill").Preload("Vendor").Preload("Product").
		Preload("OriginZone").Preload("DestZone").Preload("Mot").Preload("Uom")
	query = applyCommonContractFilters(query, filters)
	if search != "" {
		like := "%" + search + "%"
		query = query.Select("DISTINCT contract_dedicated_vars.*").
			Joins("LEFT JOIN vendors ON vendors.id = contract_dedicated_vars.vendor_id").
			Joins("LEFT JOIN mills ON mills.id = contract_dedicated_vars.mill_id").
			Joins("LEFT JOIN products ON products.id = contract_dedicated_vars.product_id").
			Joins("LEFT JOIN zones AS origin_zones ON origin_zones.id = contract_dedicated_vars.origin_zone_id").
			Joins("LEFT JOIN zones AS dest_zones ON dest_zones.id = contract_dedicated_vars.dest_zone_id").
			Joins("LEFT JOIN mots ON mots.id = contract_dedicated_vars.mot_id").
			Joins("LEFT JOIN uoms ON uoms.id = contract_dedicated_vars.uom_id").
			Where(`
				contract_dedicated_vars.spk_number ILIKE ? OR
				contract_dedicated_vars.fa_number ILIKE ? OR
				contract_dedicated_vars.area_category ILIKE ? OR
				contract_dedicated_vars.proposal_cfas ILIKE ? OR
				contract_dedicated_vars.notes ILIKE ? OR
				vendors.name ILIKE ? OR
				vendors.code ILIKE ? OR
				mills.name ILIKE ? OR
				mills.code ILIKE ? OR
				products.name ILIKE ? OR
				origin_zones.name ILIKE ? OR
				dest_zones.name ILIKE ? OR
				mots.name ILIKE ? OR
				uoms.name ILIKE ? OR
				CAST(contract_dedicated_vars.distance AS TEXT) ILIKE ? OR
				CAST(contract_dedicated_vars.payload AS TEXT) ILIKE ? OR
				CAST(contract_dedicated_vars.cost_id_r AS TEXT) ILIKE ? OR
				CAST(contract_dedicated_vars.cost_per_kg AS TEXT) ILIKE ? OR
				CAST(contract_dedicated_vars.cost_per_kgkm AS TEXT) ILIKE ?
			`,
				like, like, like, like, like,
				like, like, like, like, like, like, like, like, like,
				like, like, like, like, like,
			)
	}

	countQuery := query.Session(&gorm.Session{})
	var total int64
	if err := countQuery.Distinct("contract_dedicated_vars.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	if err := query.Order("contract_dedicated_vars.id ASC").Find(&contracts).Error; err != nil {
		return nil, 0, err
	}
	return contracts, total, nil
}

func (r *ContractRepository) GetDedicatedVarByID(id uint) (*models.ContractDedicatedVar, error) {
	var contract models.ContractDedicatedVar
	err := r.db.Preload("Mill").Preload("Vendor").Preload("Product").
		Preload("OriginZone").Preload("DestZone").Preload("Mot").Preload("Uom").
		First(&contract, id).Error
	return &contract, err
}

func (r *ContractRepository) CreateDedicatedVar(contract *models.ContractDedicatedVar) error {
	return r.db.Omit(clause.Associations).Create(contract).Error
}

func (r *ContractRepository) UpdateDedicatedVar(contract *models.ContractDedicatedVar) error {
	return updateExistingContract(r.db, &models.ContractDedicatedVar{}, contract.ID, contract)
}

func (r *ContractRepository) DeleteDedicatedVar(id uint) error {
	return deleteExistingContract(r.db, &models.ContractDedicatedVar{}, id)
}

func (r *ContractRepository) BulkCreateDedicatedVar(contracts []models.ContractDedicatedVar) error {
	return r.db.CreateInBatches(contracts, 100).Error
}

func (r *ContractRepository) FindDedicatedVarBySPK(spk string) (*models.ContractDedicatedVar, error) {
	var contract models.ContractDedicatedVar
	err := r.db.Where("spk_number = ?", spk).First(&contract).Error
	return &contract, err
}

func (r *ContractRepository) FindDedicatedVarBySPKVendorMill(spk string, vendorID, millID uint) (*models.ContractDedicatedVar, error) {
	normalized := normalizeSPK(spk)
	if normalized == "" {
		return nil, gorm.ErrRecordNotFound
	}
	var contract models.ContractDedicatedVar
	err := r.db.Where("LOWER(TRIM(spk_number)) = ? AND vendor_id = ? AND mill_id = ?", normalized, vendorID, millID).First(&contract).Error
	return &contract, err
}

// --- Oncall ---

func (r *ContractRepository) GetAllOncall(filters map[string]interface{}, search string) ([]models.ContractOncall, error) {
	var contracts []models.ContractOncall
	query := r.db.Model(&models.ContractOncall{}).Preload("Mill").Preload("Vendor").Preload("Product").
		Preload("OriginZone").Preload("DestZone").Preload("Mot").Preload("Uom")
	query = applyCommonContractFilters(query, filters)
	if search != "" {
		like := "%" + search + "%"
		query = query.Select("DISTINCT contract_oncalls.*").
			Joins("LEFT JOIN vendors ON vendors.id = contract_oncalls.vendor_id").
			Joins("LEFT JOIN mills ON mills.id = contract_oncalls.mill_id").
			Joins("LEFT JOIN products ON products.id = contract_oncalls.product_id").
			Joins("LEFT JOIN zones AS origin_zones ON origin_zones.id = contract_oncalls.origin_zone_id").
			Joins("LEFT JOIN zones AS dest_zones ON dest_zones.id = contract_oncalls.dest_zone_id").
			Joins("LEFT JOIN mots ON mots.id = contract_oncalls.mot_id").
			Joins("LEFT JOIN uoms ON uoms.id = contract_oncalls.uom_id").
			Where(`
				contract_oncalls.spk_number ILIKE ? OR
				contract_oncalls.fa_number ILIKE ? OR
				contract_oncalls.area_category ILIKE ? OR
				contract_oncalls.proposal_cfas ILIKE ? OR
				contract_oncalls.notes ILIKE ? OR
				vendors.name ILIKE ? OR
				vendors.code ILIKE ? OR
				mills.name ILIKE ? OR
				mills.code ILIKE ? OR
				products.name ILIKE ? OR
				origin_zones.name ILIKE ? OR
				dest_zones.name ILIKE ? OR
				mots.name ILIKE ? OR
				uoms.name ILIKE ? OR
				CAST(contract_oncalls.distance AS TEXT) ILIKE ? OR
				CAST(contract_oncalls.payload AS TEXT) ILIKE ? OR
				CAST(contract_oncalls.loading_cost AS TEXT) ILIKE ? OR
				CAST(contract_oncalls.unloading_cost AS TEXT) ILIKE ? OR
				CAST(contract_oncalls.cost_id_r AS TEXT) ILIKE ? OR
				CAST(contract_oncalls.cost_per_kg AS TEXT) ILIKE ? OR
				CAST(contract_oncalls.cost_per_ton AS TEXT) ILIKE ? OR
				CAST(contract_oncalls.running_cost_id_r AS TEXT) ILIKE ? OR
				CAST(contract_oncalls.running_cost_usd AS TEXT) ILIKE ?
			`,
				like, like, like, like, like,
				like, like, like, like, like, like, like, like, like,
				like, like, like, like, like, like, like, like, like,
			)
	}
	err := query.Order("contract_oncalls.id ASC").Find(&contracts).Error
	return contracts, err
}

func (r *ContractRepository) GetOncallPage(filters map[string]interface{}, search string, limit, offset int) ([]models.ContractOncall, int64, error) {
	var contracts []models.ContractOncall
	query := r.db.Model(&models.ContractOncall{}).Preload("Mill").Preload("Vendor").Preload("Product").
		Preload("OriginZone").Preload("DestZone").Preload("Mot").Preload("Uom")
	query = applyCommonContractFilters(query, filters)
	if search != "" {
		like := "%" + search + "%"
		query = query.Select("DISTINCT contract_oncalls.*").
			Joins("LEFT JOIN vendors ON vendors.id = contract_oncalls.vendor_id").
			Joins("LEFT JOIN mills ON mills.id = contract_oncalls.mill_id").
			Joins("LEFT JOIN products ON products.id = contract_oncalls.product_id").
			Joins("LEFT JOIN zones AS origin_zones ON origin_zones.id = contract_oncalls.origin_zone_id").
			Joins("LEFT JOIN zones AS dest_zones ON dest_zones.id = contract_oncalls.dest_zone_id").
			Joins("LEFT JOIN mots ON mots.id = contract_oncalls.mot_id").
			Joins("LEFT JOIN uoms ON uoms.id = contract_oncalls.uom_id").
			Where(`
				contract_oncalls.spk_number ILIKE ? OR
				contract_oncalls.fa_number ILIKE ? OR
				contract_oncalls.area_category ILIKE ? OR
				contract_oncalls.proposal_cfas ILIKE ? OR
				contract_oncalls.notes ILIKE ? OR
				vendors.name ILIKE ? OR
				vendors.code ILIKE ? OR
				mills.name ILIKE ? OR
				mills.code ILIKE ? OR
				products.name ILIKE ? OR
				origin_zones.name ILIKE ? OR
				dest_zones.name ILIKE ? OR
				mots.name ILIKE ? OR
				uoms.name ILIKE ? OR
				CAST(contract_oncalls.distance AS TEXT) ILIKE ? OR
				CAST(contract_oncalls.payload AS TEXT) ILIKE ? OR
				CAST(contract_oncalls.loading_cost AS TEXT) ILIKE ? OR
				CAST(contract_oncalls.unloading_cost AS TEXT) ILIKE ? OR
				CAST(contract_oncalls.cost_id_r AS TEXT) ILIKE ? OR
				CAST(contract_oncalls.cost_per_kg AS TEXT) ILIKE ? OR
				CAST(contract_oncalls.cost_per_ton AS TEXT) ILIKE ? OR
				CAST(contract_oncalls.running_cost_id_r AS TEXT) ILIKE ? OR
				CAST(contract_oncalls.running_cost_usd AS TEXT) ILIKE ?
			`,
				like, like, like, like, like,
				like, like, like, like, like, like, like, like, like,
				like, like, like, like, like, like, like, like, like,
			)
	}

	countQuery := query.Session(&gorm.Session{})
	var total int64
	if err := countQuery.Distinct("contract_oncalls.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	if err := query.Order("contract_oncalls.id ASC").Find(&contracts).Error; err != nil {
		return nil, 0, err
	}
	return contracts, total, nil
}

func (r *ContractRepository) GetOncallByID(id uint) (*models.ContractOncall, error) {
	var contract models.ContractOncall
	err := r.db.Preload("Mill").Preload("Vendor").Preload("Product").
		Preload("OriginZone").Preload("DestZone").Preload("Mot").Preload("Uom").
		First(&contract, id).Error
	return &contract, err
}

func (r *ContractRepository) CreateOncall(contract *models.ContractOncall) error {
	return r.db.Omit(clause.Associations).Create(contract).Error
}

func (r *ContractRepository) UpdateOncall(contract *models.ContractOncall) error {
	return updateExistingContract(r.db, &models.ContractOncall{}, contract.ID, contract)
}

func (r *ContractRepository) DeleteOncall(id uint) error {
	return deleteExistingContract(r.db, &models.ContractOncall{}, id)
}

func (r *ContractRepository) BulkCreateOncall(contracts []models.ContractOncall) error {
	return r.db.CreateInBatches(contracts, 100).Error
}

func (r *ContractRepository) FindOncallBySPK(spk string) (*models.ContractOncall, error) {
	var contract models.ContractOncall
	err := r.db.Where("spk_number = ?", spk).First(&contract).Error
	return &contract, err
}

func (r *ContractRepository) FindOncallBySPKVendorMill(spk string, vendorID, millID uint) (*models.ContractOncall, error) {
	normalized := normalizeSPK(spk)
	if normalized == "" {
		return nil, gorm.ErrRecordNotFound
	}
	var contract models.ContractOncall
	err := r.db.Where("LOWER(TRIM(spk_number)) = ? AND vendor_id = ? AND mill_id = ?", normalized, vendorID, millID).First(&contract).Error
	return &contract, err
}

// --- Map-based Updates (partial update, no associations) ---

func (r *ContractRepository) UpdateDedicatedFixMap(id uint, updates map[string]interface{}) error {
	return updateExistingContractMap(r.db, &models.ContractDedicatedFix{}, id, updates)
}

func (r *ContractRepository) UpdateDedicatedVarMap(id uint, updates map[string]interface{}) error {
	return updateExistingContractMap(r.db, &models.ContractDedicatedVar{}, id, updates)
}

func (r *ContractRepository) UpdateOncallMap(id uint, updates map[string]interface{}) error {
	return updateExistingContractMap(r.db, &models.ContractOncall{}, id, updates)
}
