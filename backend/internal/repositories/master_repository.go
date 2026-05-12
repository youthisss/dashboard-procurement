package repositories

import (
	"crypto/sha1"
	"encoding/hex"
	"rygell-dashboard/internal/models"
	"strings"
	"unicode"

	"gorm.io/gorm"
)

// MasterRepository handles CRUD for all master data entities.
type MasterRepository struct {
	db *gorm.DB
}

// NewMasterRepository creates a new MasterRepository.
func NewMasterRepository(db *gorm.DB) *MasterRepository {
	return &MasterRepository{db: db}
}

// DB exposes the current GORM handle for services that need to create
// transaction-scoped repositories.
func (r *MasterRepository) DB() *gorm.DB {
	return r.db
}

func updateExisting(db *gorm.DB, model interface{}, id uint, columns []string, values interface{}) error {
	if err := db.First(model, id).Error; err != nil {
		return err
	}
	return db.Model(model).Select(columns).Updates(values).Error
}

func deleteExisting(db *gorm.DB, model interface{}, id uint) error {
	result := db.Delete(model, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// --- Mill ---

func (r *MasterRepository) GetAllMills(search string) ([]models.Mill, error) {
	var mills []models.Mill
	query := r.db
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ?", like, like)
	}
	err := query.Order("id ASC").Find(&mills).Error
	return mills, err
}

// SearchMillsLimited returns at most `limit` mills matching the search query.
func (r *MasterRepository) SearchMillsLimited(search string, limit int) ([]models.Mill, error) {
	var mills []models.Mill
	query := r.db
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ?", like, like)
	}
	err := query.Order("id ASC").Limit(limit).Find(&mills).Error
	return mills, err
}

func (r *MasterRepository) GetMillByID(id uint) (*models.Mill, error) {
	var mill models.Mill
	err := r.db.First(&mill, id).Error
	return &mill, err
}

func (r *MasterRepository) CreateMill(mill *models.Mill) error {
	return r.db.Create(mill).Error
}

func (r *MasterRepository) UpdateMill(mill *models.Mill) error {
	return updateExisting(r.db, &models.Mill{}, mill.ID, []string{"code", "name"}, mill)
}

func (r *MasterRepository) DeleteMill(id uint) error {
	return deleteExisting(r.db, &models.Mill{}, id)
}

// --- Vendor ---

func (r *MasterRepository) GetAllVendors(search string) ([]models.Vendor, error) {
	var vendors []models.Vendor
	query := r.db
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ? OR email ILIKE ? OR contact_person ILIKE ? OR phone ILIKE ? OR address ILIKE ?", like, like, like, like, like, like)
	}
	err := query.Order("id ASC").Find(&vendors).Error
	return vendors, err
}

// SearchVendorsLimited returns at most `limit` vendors matching the search query.
func (r *MasterRepository) SearchVendorsLimited(search string, limit int) ([]models.Vendor, error) {
	var vendors []models.Vendor
	query := r.db
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ? OR email ILIKE ? OR contact_person ILIKE ? OR phone ILIKE ? OR address ILIKE ?", like, like, like, like, like, like)
	}
	err := query.Order("id ASC").Limit(limit).Find(&vendors).Error
	return vendors, err
}

func (r *MasterRepository) GetVendorByID(id uint) (*models.Vendor, error) {
	var vendor models.Vendor
	err := r.db.First(&vendor, id).Error
	return &vendor, err
}

func (r *MasterRepository) CreateVendor(vendor *models.Vendor) error {
	return r.db.Create(vendor).Error
}

func (r *MasterRepository) UpdateVendor(vendor *models.Vendor) error {
	return updateExisting(
		r.db,
		&models.Vendor{},
		vendor.ID,
		[]string{"code", "name", "tax_id", "status", "address", "contact_person", "email", "phone"},
		vendor,
	)
}

func (r *MasterRepository) DeleteVendor(id uint) error {
	return deleteExisting(r.db, &models.Vendor{}, id)
}

// --- Product ---

func (r *MasterRepository) GetAllProducts(search string) ([]models.Product, error) {
	var products []models.Product
	query := r.db
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name ILIKE ?", like)
	}
	err := query.Order("id ASC").Find(&products).Error
	return products, err
}

func (r *MasterRepository) GetProductByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error
	return &product, err
}

func (r *MasterRepository) CreateProduct(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *MasterRepository) UpdateProduct(product *models.Product) error {
	return updateExisting(r.db, &models.Product{}, product.ID, []string{"name"}, product)
}

func (r *MasterRepository) DeleteProduct(id uint) error {
	return deleteExisting(r.db, &models.Product{}, id)
}

// --- Zone ---

func (r *MasterRepository) GetAllZones(search string) ([]models.Zone, error) {
	var zones []models.Zone
	query := r.db
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name ILIKE ? OR type ILIKE ?", like, like)
	}
	err := query.Order("id ASC").Find(&zones).Error
	return zones, err
}

// SearchZonesLimited returns at most `limit` zones matching the search query.
func (r *MasterRepository) SearchZonesLimited(search string, limit int) ([]models.Zone, error) {
	var zones []models.Zone
	query := r.db
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name ILIKE ? OR type ILIKE ?", like, like)
	}
	err := query.Order("id ASC").Limit(limit).Find(&zones).Error
	return zones, err
}

func (r *MasterRepository) GetZoneByID(id uint) (*models.Zone, error) {
	var zone models.Zone
	err := r.db.First(&zone, id).Error
	return &zone, err
}

func (r *MasterRepository) CreateZone(zone *models.Zone) error {
	return r.db.Create(zone).Error
}

func (r *MasterRepository) UpdateZone(zone *models.Zone) error {
	return updateExisting(r.db, &models.Zone{}, zone.ID, []string{"name", "type"}, zone)
}

func (r *MasterRepository) DeleteZone(id uint) error {
	return deleteExisting(r.db, &models.Zone{}, id)
}

// --- MOT ---

func (r *MasterRepository) GetAllMots(search string) ([]models.Mot, error) {
	var mots []models.Mot
	query := r.db
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name ILIKE ?", like)
	}
	err := query.Order("id ASC").Find(&mots).Error
	return mots, err
}

func (r *MasterRepository) GetMotByID(id uint) (*models.Mot, error) {
	var mot models.Mot
	err := r.db.First(&mot, id).Error
	return &mot, err
}

func (r *MasterRepository) CreateMot(mot *models.Mot) error {
	return r.db.Create(mot).Error
}

func (r *MasterRepository) UpdateMot(mot *models.Mot) error {
	return updateExisting(r.db, &models.Mot{}, mot.ID, []string{"name"}, mot)
}

func (r *MasterRepository) DeleteMot(id uint) error {
	return deleteExisting(r.db, &models.Mot{}, id)
}

// --- UOM ---

func (r *MasterRepository) GetAllUoms(search string) ([]models.Uom, error) {
	var uoms []models.Uom
	query := r.db
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name ILIKE ?", like)
	}
	err := query.Order("id ASC").Find(&uoms).Error
	return uoms, err
}

func (r *MasterRepository) GetUomByID(id uint) (*models.Uom, error) {
	var uom models.Uom
	err := r.db.First(&uom, id).Error
	return &uom, err
}

func (r *MasterRepository) CreateUom(uom *models.Uom) error {
	return r.db.Create(uom).Error
}

func (r *MasterRepository) UpdateUom(uom *models.Uom) error {
	return updateExisting(r.db, &models.Uom{}, uom.ID, []string{"name"}, uom)
}

func (r *MasterRepository) DeleteUom(id uint) error {
	return deleteExisting(r.db, &models.Uom{}, id)
}

// --- Bulk Operations ---

func (r *MasterRepository) BulkCreateMills(mills []models.Mill) error {
	return r.db.CreateInBatches(mills, 100).Error
}

func (r *MasterRepository) BulkCreateVendors(vendors []models.Vendor) error {
	return r.db.CreateInBatches(vendors, 100).Error
}

func (r *MasterRepository) BulkCreateProducts(products []models.Product) error {
	return r.db.CreateInBatches(products, 100).Error
}

func (r *MasterRepository) BulkCreateZones(zones []models.Zone) error {
	return r.db.CreateInBatches(zones, 100).Error
}

// --- Find or Create (Upsert by Name) ---

// FindOrCreateVendorByName finds a vendor by name, or creates one if not found.
func (r *MasterRepository) FindOrCreateVendorByName(name string) (*models.Vendor, error) {
	var vendor models.Vendor
	err := r.db.Where("name = ?", name).First(&vendor).Error
	if err == nil {
		return &vendor, nil
	}
	vendor = models.Vendor{Name: name, Code: importCodeFromName("V", name)}
	if err := r.db.Create(&vendor).Error; err != nil {
		return nil, err
	}
	return &vendor, nil
}

// FindOrCreateMillByName finds a mill by name, or creates one if not found.
func (r *MasterRepository) FindOrCreateMillByName(name string) (*models.Mill, error) {
	var mill models.Mill
	err := r.db.Where("name = ?", name).First(&mill).Error
	if err == nil {
		return &mill, nil
	}
	mill = models.Mill{Name: name, Code: importCodeFromName("M", name)}
	if err := r.db.Create(&mill).Error; err != nil {
		return nil, err
	}
	return &mill, nil
}

// FindOrCreateProductByName finds a product by name, or creates one if not found.
func (r *MasterRepository) FindOrCreateProductByName(name string) (*models.Product, error) {
	var product models.Product
	err := r.db.Where("name = ?", name).First(&product).Error
	if err == nil {
		return &product, nil
	}
	product = models.Product{Name: name}
	if err := r.db.Create(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// FindOrCreateZoneByName finds a zone by name, or creates one if not found.
func (r *MasterRepository) FindOrCreateZoneByName(name string) (*models.Zone, error) {
	var zone models.Zone
	err := r.db.Where("name = ?", name).First(&zone).Error
	if err == nil {
		return &zone, nil
	}
	zone = models.Zone{Name: name}
	if err := r.db.Create(&zone).Error; err != nil {
		return nil, err
	}
	return &zone, nil
}

// FindOrCreateMotByName finds a MOT by name, or creates one if not found.
func (r *MasterRepository) FindOrCreateMotByName(name string) (*models.Mot, error) {
	var mot models.Mot
	err := r.db.Where("name = ?", name).First(&mot).Error
	if err == nil {
		return &mot, nil
	}
	mot = models.Mot{Name: name}
	if err := r.db.Create(&mot).Error; err != nil {
		return nil, err
	}
	return &mot, nil
}

// FindOrCreateUomByName finds a UOM by name, or creates one if not found.
func (r *MasterRepository) FindOrCreateUomByName(name string) (*models.Uom, error) {
	var uom models.Uom
	err := r.db.Where("name = ?", name).First(&uom).Error
	if err == nil {
		return &uom, nil
	}
	uom = models.Uom{Name: name}
	if err := r.db.Create(&uom).Error; err != nil {
		return nil, err
	}
	return &uom, nil
}

func importCodeFromName(prefix, name string) string {
	const maxCodeLength = 50

	prefix = strings.ToUpper(strings.TrimSpace(prefix))
	if prefix == "" {
		prefix = "X"
	}

	base := normalizeCodeBase(name)
	if base == "" {
		base = "ITEM"
	}

	sum := sha1.Sum([]byte(strings.ToLower(strings.TrimSpace(name))))
	suffix := strings.ToUpper(hex.EncodeToString(sum[:4]))
	reserved := len(prefix) + 2 + len(suffix)
	maxBaseLength := maxCodeLength - reserved
	if maxBaseLength < 1 {
		maxBaseLength = 1
	}
	if len(base) > maxBaseLength {
		base = strings.Trim(base[:maxBaseLength], "-")
		if base == "" {
			base = "ITEM"
		}
	}

	return prefix + "-" + base + "-" + suffix
}

func normalizeCodeBase(name string) string {
	var builder strings.Builder
	lastDash := false
	for _, r := range strings.TrimSpace(name) {
		r = unicode.ToUpper(r)
		switch {
		case (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9'):
			builder.WriteRune(r)
			lastDash = false
		default:
			if !lastDash && builder.Len() > 0 {
				builder.WriteByte('-')
				lastDash = true
			}
		}
	}
	return strings.Trim(builder.String(), "-")
}
