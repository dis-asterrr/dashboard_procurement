package repositories

import (
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

// --- Dedicated Fix ---

func (r *ContractRepository) GetAllDedicatedFix(filters map[string]interface{}, search string) ([]models.ContractDedicatedFix, error) {
	var contracts []models.ContractDedicatedFix
	query := r.db.Preload("Mill").Preload("Vendor").Preload("Product")
	if v, ok := filters["vendor_id"]; ok && v != "" {
		query = query.Where("vendor_id = ?", v)
	}
	if v, ok := filters["mill_id"]; ok && v != "" {
		query = query.Where("mill_id = ?", v)
	}
	if search != "" {
		like := "%" + search + "%"
		query = query.Select("DISTINCT contract_dedicated_fixes.*").
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
	err := query.Order("contract_dedicated_fixes.id ASC").Find(&contracts).Error
	return contracts, err
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
	return r.db.Omit(clause.Associations).Save(contract).Error
}

func (r *ContractRepository) DeleteDedicatedFix(id uint) error {
	return r.db.Delete(&models.ContractDedicatedFix{}, id).Error
}

func (r *ContractRepository) BulkCreateDedicatedFix(contracts []models.ContractDedicatedFix) error {
	return r.db.CreateInBatches(contracts, 100).Error
}

func (r *ContractRepository) FindDedicatedFixBySPK(spk string) (*models.ContractDedicatedFix, error) {
	var contract models.ContractDedicatedFix
	err := r.db.Where("spk_number = ?", spk).First(&contract).Error
	return &contract, err
}

// --- Dedicated Var ---

func (r *ContractRepository) GetAllDedicatedVar(filters map[string]interface{}, search string) ([]models.ContractDedicatedVar, error) {
	var contracts []models.ContractDedicatedVar
	query := r.db.Preload("Mill").Preload("Vendor").Preload("Product").
		Preload("OriginZone").Preload("DestZone").Preload("Mot").Preload("Uom")
	if v, ok := filters["vendor_id"]; ok && v != "" {
		query = query.Where("vendor_id = ?", v)
	}
	if v, ok := filters["mill_id"]; ok && v != "" {
		query = query.Where("mill_id = ?", v)
	}
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
				CAST(contract_dedicated_vars.cost_idr AS TEXT) ILIKE ? OR
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
	return r.db.Omit(clause.Associations).Save(contract).Error
}

func (r *ContractRepository) DeleteDedicatedVar(id uint) error {
	return r.db.Delete(&models.ContractDedicatedVar{}, id).Error
}

func (r *ContractRepository) BulkCreateDedicatedVar(contracts []models.ContractDedicatedVar) error {
	return r.db.CreateInBatches(contracts, 100).Error
}

func (r *ContractRepository) FindDedicatedVarBySPK(spk string) (*models.ContractDedicatedVar, error) {
	var contract models.ContractDedicatedVar
	err := r.db.Where("spk_number = ?", spk).First(&contract).Error
	return &contract, err
}

// --- Oncall ---

func (r *ContractRepository) GetAllOncall(filters map[string]interface{}, search string) ([]models.ContractOncall, error) {
	var contracts []models.ContractOncall
	query := r.db.Preload("Mill").Preload("Vendor").Preload("Product").
		Preload("OriginZone").Preload("DestZone").Preload("Mot").Preload("Uom")
	if v, ok := filters["vendor_id"]; ok && v != "" {
		query = query.Where("vendor_id = ?", v)
	}
	if v, ok := filters["mill_id"]; ok && v != "" {
		query = query.Where("mill_id = ?", v)
	}
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
				CAST(contract_oncalls.cost_idr AS TEXT) ILIKE ? OR
				CAST(contract_oncalls.cost_per_kg AS TEXT) ILIKE ? OR
				CAST(contract_oncalls.cost_per_ton AS TEXT) ILIKE ? OR
				CAST(contract_oncalls.running_cost_idr AS TEXT) ILIKE ? OR
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
	return r.db.Omit(clause.Associations).Save(contract).Error
}

func (r *ContractRepository) DeleteOncall(id uint) error {
	return r.db.Delete(&models.ContractOncall{}, id).Error
}

func (r *ContractRepository) BulkCreateOncall(contracts []models.ContractOncall) error {
	return r.db.CreateInBatches(contracts, 100).Error
}

func (r *ContractRepository) FindOncallBySPK(spk string) (*models.ContractOncall, error) {
	var contract models.ContractOncall
	err := r.db.Where("spk_number = ?", spk).First(&contract).Error
	return &contract, err
}

// --- Map-based Updates (partial update, no associations) ---

func (r *ContractRepository) UpdateDedicatedFixMap(id uint, updates map[string]interface{}) error {
	return r.db.Model(&models.ContractDedicatedFix{}).Where("id = ?", id).Updates(updates).Error
}

func (r *ContractRepository) UpdateDedicatedVarMap(id uint, updates map[string]interface{}) error {
	return r.db.Model(&models.ContractDedicatedVar{}).Where("id = ?", id).Updates(updates).Error
}

func (r *ContractRepository) UpdateOncallMap(id uint, updates map[string]interface{}) error {
	return r.db.Model(&models.ContractOncall{}).Where("id = ?", id).Updates(updates).Error
}
