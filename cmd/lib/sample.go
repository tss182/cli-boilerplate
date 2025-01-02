package lib

var SampleEntity = `
package entity

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"{{.GoModName}}/shared"
	"{{.GoModName}}/shared/pagination"
)

var Name string = "{{.Domain}}"

type (
	// Interface method repository
	{{.DomainPackage}}RepositoryInterface interface {
		Get(ctx context.Context, dto *pagination.Request) ({{.DomainPackage}}Response, error)
		GetByID(ctx context.Context, code *string) ({{.DomainPackage}}Model, error)
		Create(ctx context.Context, dto *{{.DomainPackage}}Model) error
		Update(ctx context.Context, dto *{{.DomainPackage}}Model) error
		Delete(ctx context.Context, code *string, userLog string) error
	}

	// Interface method feature
	{{.DomainPackage}}FeatureInterface interface {
		Get(ctx context.Context, dto *pagination.Request) ({{.DomainPackage}}Response, error)
		GetByID(ctx context.Context, code *string) ({{.DomainPackage}}Model, error)
		Create(ctx context.Context, dto *{{.DomainPackage}}Model) error
		Update(ctx context.Context, dto *{{.DomainPackage}}Model) error
		Delete(ctx context.Context, code *string, userLog string) error
	}

	// Interface method handler
	{{.DomainPackage}}HandlerInterface interface {
		Get(c *fiber.Ctx) error
		GetByID(c *fiber.Ctx) error
		Create(c *fiber.Ctx) error
		Update(c *fiber.Ctx) error
		Delete(c *fiber.Ctx) error
	}

	// Model
	{{.DomainPackage}}CodeModel struct {
		Code string ` + "`" + `db:"code" json:"code"` + "`" + `
	}


	{{.DomainPackage}}HeaderModel struct {
		Name string ` + "`" + `validate:"required,max=50" db:"name" json:"name"` + "`" + `
	}

	{{.DomainPackage}}Model struct {
		{{.DomainPackage}}CodeModel
		{{.DomainPackage}}HeaderModel
		shared.LogModel
	}

	// Request
	{{.DomainPackage}}Request struct {
		{{.DomainPackage}}CodeModel
		{{.DomainPackage}}HeaderModel
	}

	// Response
	{{.DomainPackage}}Response struct {
		Data []{{.DomainPackage}}Model ` + "`" + `json:"rows"` + "`" + `
		pagination.Response
	}
)

`

var SampleRepository = `
package repository

import (
	"context"
	"{{.GoModName}}/config"
	"{{.GoModName}}/{{.Folder}}/entity"
	"{{.GoModName}}/infrastructure/database"
	"{{.GoModName}}/shared"
	"{{.GoModName}}/shared/constant"
	"{{.GoModName}}/shared/pagination"
	"time"
	"github.com/huandu/go-sqlbuilder"
)

type {{.DomainPackage}}RepositoryInterface interface {
	entity.{{.DomainPackage}}RepositoryInterface
}

type {{.DomainPackageLocal}}Repository struct {
	db     *database.Database
	config *config.Config
}

const tableName = "name_table"

func New(dbMySQL *database.Database, cfg *config.Config) {{.DomainPackage}}RepositoryInterface {
	return &{{.DomainPackageLocal}}Repository{
		db:     dbMySQL,
		config: cfg,
	}
}

func (r *{{.DomainPackageLocal}}Repository) Get(ctx context.Context, dto *pagination.Request) (res entity.{{.DomainPackage}}Response, e error) {
	var data []entity.{{.DomainPackage}}Model
	var fields = []string{"code", "name"}

	sb := sqlbuilder.NewSelectBuilder()
	query := sb.Select(
		"IFNULL(code, '') AS code",
		"IFNULL(name, '') AS name",
		created_at,
		shared.SQLSelectFullName(tableName+".created_by", "created_by"),
		modified_at,
		shared.SQLSelectFullName(tableName+".modified_by", "modified_by"),
	).
		From(tableName).
		Where(sb.IsNull("deleted_at")).
		OrderBy("code ASC")

	if dto.Search != "" {
		var orConditions []string
		for _, field := range fields {
			orConditions = append(orConditions, sb.Like(field, "%"+dto.Search+"%"))
		}
		sql = sql.Where(sb.Or(orConditions...))
	}

	queryPage, pMeta, err := pagination.New(query, r.config, dto)
	if err != nil {
		e = err
		return
	}

	err = r.db.DB.SelectContext(ctx, &data, queryPage.Raw, queryPage.Args...)
	if err != nil {
		e = err
		return
	}

	err = r.db.DB.GetContext(ctx, &pMeta.TotalRows, queryPage.Count, queryPage.Args...)
	if err != nil {
		e = err
		return
	}

	queryPage.SetTotal(pMeta.TotalRows, pMeta.Limit, &pMeta.TotalPages)
	res.Data = data
	res.PaginationMeta = pMeta
	return
}

func (r *{{.DomainPackageLocal}}Repository) GetByID(ctx context.Context, code *string) (res entity.{{.DomainPackage}}Model, e error) {
	sb := sqlbuilder.NewSelectBuilder()
	query, args := sb.Select(
		"IFNULL(code, '') AS code",
		"IFNULL(name, '') AS name",
		created_at,
		shared.SQLSelectFullName(tableName+".created_by", "created_by"),
		modified_at,
		shared.SQLSelectFullName(tableName+".modified_by", "modified_by"),
	).
		From(tableName).
		Where(
			sb.Equal("code", code),
			sb.IsNull("deleted_at"),
		).
		OrderBy("created_at DESC").
		Build()

	err := r.db.DB.GetContext(ctx, &res, query, args...)
	if err != nil {
		e = err
		return
	}

	return
}

func (r *{{.DomainPackageLocal}}Repository) Create(ctx context.Context, dto *entity.{{.DomainPackage}}Model) error {
	ib := sqlbuilder.NewInsertBuilder()
	query, args := ib.InsertInto(tableName).
		Cols("code", "name", "created_by", "created_at").
		Values(dto.Code, dto.Name, dto.CreatedBy, dto.CreatedAt).
		Build()

	_, err := r.db.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *{{.DomainPackageLocal}}Repository) Update(ctx context.Context, dto *entity.{{.DomainPackage}}Model) error {
	ub := sqlbuilder.NewUpdateBuilder()
	query, args := ub.Update(tableName).
		Set(
			ub.Assign("name", dto.Name),
			ub.Assign("modified_by", dto.ModifiedBy),
			ub.Assign("modified_at", dto.ModifiedAt),
		).
		Where(
			ub.Equal("code", dto.Code),
			ub.IsNull("deleted_at"),
		).
		Build()

	_, err := r.db.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *{{.DomainPackageLocal}}Repository) Delete(ctx context.Context, code *string, userLog string) error {
	ub := sqlbuilder.NewUpdateBuilder()
	query, args := ub.Update(tableName).
		Set(
			ub.Assign("deleted_at", time.Now().UTC()),
			ub.Assign("deleted_by", userLog),
		).
		Where(
			ub.Equal("code", code),
			ub.IsNull("deleted_at"),
		).
		Build()

	_, err := r.db.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
`

var SampleFeature = `
package feature

import (
	"context"
	"database/sql"
	"{{.GoModName}}/{{.Folder}}/entity"
	"{{.GoModName}}/shared/pagination"
	"errors"
	"fmt"
)

type {{.DomainPackage}}FeatureInterface interface {
	entity.{{.DomainPackage}}FeatureInterface
}

type {{.DomainPackageLocal}}Feature struct {
	repo entity.{{.DomainPackage}}RepositoryInterface
}

func New(repo entity.{{.DomainPackage}}RepositoryInterface) {{.DomainPackage}}FeatureInterface {
	return &{{.DomainPackageLocal}}Feature{
		repo: repo,
	}
}

func (r *{{.DomainPackageLocal}}Feature) Get(ctx context.Context, dto *pagination.Request) (entity.{{.DomainPackage}}Response, error) {
	tz := ctx.Value(constant.CTX_TIMEZONE).(string)
	timezone := shared.GetTimeZone(tz)
	res, err := r.repo.Get(ctx, dto)
	if err != nil {
		return res, err
	}
	for i, v := range res.{{.DomainPackage}} {
		res.{{.DomainPackage}}[i].CreatedAt = v.CreatedAt.In(timezone)
		res.{{.DomainPackage}}[i].ModifiedAt = shared.SetTimeZone(v.ModifiedAt, timezone)
	}
	return res, nil
}

func (r *{{.DomainPackageLocal}}Feature) GetByID(ctx context.Context, code *string) (entity.{{.DomainPackage}}Model, error) {
	tz := ctx.Value(constant.CTX_TIMEZONE).(string)
	timezone := shared.GetTimeZone(tz)
	res, err := r.repo.GetByID(ctx, code)
	if err != nil {
		return res, err
	}
	res.CreatedAt = res.CreatedAt.In(timezone)
	res.ModifiedAt = shared.SetTimeZone(res.ModifiedAt, timezone)
	return res, nil
}

func (r *{{.DomainPackageLocal}}Feature) Create(ctx context.Context, dto *entity.{{.DomainPackage}}Model) error {
	data, err := r.GetByID(ctx, &dto.Code)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if data.Code != "" {
		return fmt.Errorf("400:Product %s already exists", dto.Code)
	}

	return r.repo.Create(ctx, dto)
}

func (r *{{.DomainPackageLocal}}Feature) Update(ctx context.Context, dto *entity.{{.DomainPackage}}Model) error {
	_, err := r.GetByID(ctx, &dto.Code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("404:Data %s not found", dto.Code)
		} else {
			return err
		}
	}

	return r.repo.Update(ctx, dto)
}

func (r *{{.DomainPackageLocal}}Feature) Delete(ctx context.Context, code *string, userLog string) error {
	_, err := r.GetByID(ctx, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("404:Data %s not found", *code)
		} else {
			return err
		}
	}

	return r.repo.Delete(ctx, code, userLog)
}
`

var SampleHandler = `
package {{.DomainPackageLocal}}

import (
	"{{.GoModName}}/{{.Folder}}/entity"
	"{{.GoModName}}/shared"
	"{{.GoModName}}/shared/constant"
	"{{.GoModName}}/shared/message"
	"{{.GoModName}}/shared/pagination"
	"{{.GoModName}}/shared/response"
	"{{.GoModName}}/shared/validate"
	"fmt"
	"time"
	"github.com/gofiber/fiber/v2"
)

type {{.DomainPackage}}HandlerInterface interface {
	entity.{{.DomainPackage}}HandlerInterface
}

type {{.DomainPackageLocal}}Handler struct {
	feat entity.{{.DomainPackage}}FeatureInterface
}

func New(feat entity.{{.DomainPackage}}FeatureInterface) {{.DomainPackage}}HandlerInterface {
	return &{{.DomainPackageLocal}}Handler{
		feat: feat,
	}
}

// Get godoc
// @Summary {{.Domain}} List
// @Description Get all of {{.Domain}} data
// @Tags Master {{.Domain}}
// @Accept json
// @Produce json
// @Param   page	query    int     false  "Page"
// @Param   limit	query    int     false  "Limit"
// @Param   search	query    string     false  "Search"
// @Param   sortby	query    string     false  "Sort By"
// @Param   sorttype	query    string     false  "Sort Type (ASC/DESC)"
// @Success 200 {object} response.Response{data=pagination.Response{rows=[]entity.{{.DomainPackage}}Model}}
// @Failure 400,404,500 {object} response.Response
// @Router /master/data [get]
func (h *{{.DomainPackageLocal}}Handler) Get(c *fiber.Ctx) error {
	var dto pagination.Request
	err := c.QueryParser(&dto)
	shared.Exception(err, fmt.Sprintf(message.MSG_PAYLOAD, entity.Name))

	ctx, cancel := shared.CreateContextWithTimeoutAndValue(c)
	defer cancel()

	res, err := h.feat.Get(ctx, &dto)
	shared.Exception(err, entity.Name)

	if res.Data == nil {
		res.Data = []entity.{{.DomainPackage}}Model{}
	}

	return response.OK(c, fmt.Sprintf(message.MSG_LIST, entity.Name), res)
}

// GetByID godoc
// @Summary {{.Domain}} Data
// @Description Get {{.Domain}} data by code
// @Tags Master {{.Domain}}
// @Accept json
// @Produce json
// @Param   code	path    string     true  "Code"
// @Success 200 {object} entity.{{.DomainPackage}}Model
// @Failure 400,404,500 {object} response.Response
// @Router /master/data/{code} [get]
func (h *{{.DomainPackageLocal}}Handler) GetByID(c *fiber.Ctx) error {
	dto := entity.{{.DomainPackage}}CodeModel{
		Code: c.Params("code"),
	}

	errs := validate.ValidateStruct(dto)
	if len(errs) > 0 {
		return response.BadRequest(c, message.DATA_VALIDATION, errs)
	}

	ctx, cancel := shared.CreateContextWithTimeoutAndValue(c)
	defer cancel()

	res, err := h.feat.GetByID(ctx, &dto.Code)
	shared.Exception(err, entity.Name)

	return response.OK(c, fmt.Sprintf(message.MSG_DATA, entity.Name), res)
}

// Create godoc
// @Summary Create {{.Domain}}
// @Description Create new {{.Domain}} data
// @Tags Master {{.Domain}}
// @Accept json
// @Produce json
// @Param payload body entity.{{.DomainPackage}}Request true  "Payload"
// @Success 200 {object} response.Response
// @Failure 201,400,404,500 {object} response.Response
// @Router /master/data [post]
func (h *{{.DomainPackageLocal}}Handler) Create(c *fiber.Ctx) error {
	var dto entity.{{.DomainPackage}}Model
	err := c.BodyParser(&dto)
	shared.Exception(err, fmt.Sprintf(message.MSG_PAYLOAD, entity.Name))

	errs := validate.ValidateStruct(dto)
	if len(errs) > 0 {
		return response.BadRequest(c, message.DATA_VALIDATION, errs)
	}

	dto.CreatedAt = time.Now().UTC()
	dto.CreatedBy = c.Locals(constant.REQ_USERNAME).(string)

	ctx, cancel := shared.CreateContextWithTimeoutAndValue(c)
	defer cancel()

	err = h.feat.Create(ctx, &dto)
	shared.Exception(err, entity.Name)

	return response.Created(c, fmt.Sprintf(message.MSG_CREATED, entity.Name))
}

// Update godoc
// @Summary Update {{.Domain}}
// @Description Update existing {{.Domain}} data
// @Tags Master {{.Domain}}
// @Accept json
// @Produce json
// @Param payload body entity.{{.DomainPackage}}Request true  "Payload"
// @Success 200 {object} response.Response
// @Failure 400,404,500 {object} response.Response
// @Router /master/data [put]
func (h *{{.DomainPackageLocal}}Handler) Update(c *fiber.Ctx) error {
	var dto entity.{{.DomainPackage}}Model
	err := c.BodyParser(&dto)
	shared.Exception(err, fmt.Sprintf(message.MSG_PAYLOAD, entity.Name))

	errs := validate.ValidateStruct(dto)
	if len(errs) > 0 {
		return response.BadRequest(c, message.DATA_VALIDATION, errs)
	}

	now := time.Now().UTC()
	dto.ModifiedAt = &now
	dto.ModifiedBy = c.Locals(constant.REQ_USERNAME).(string)

	ctx, cancel := shared.CreateContextWithTimeoutAndValue(c)
	defer cancel()

	err = h.feat.Update(ctx, &dto)
	shared.Exception(err, entity.Name)

	return response.OK(c, fmt.Sprintf(message.MSG_UPDATED, entity.Name), nil)
}

// Delete godoc
// @Summary Delete {{.Domain}}
// @Description Delete existing {{.Domain}} data
// @Tags Master {{.Domain}}
// @Accept json
// @Produce json
// @Param payload body entity.{{.DomainPackageName}}CodeModel true  "Payload (code)"
// @Success 200 {object} response.Response
// @Failure 400,404,500 {object} response.Response
// @Router /master/data [delete]
func (h *{{.DomainPackageLocal}}Handler) Delete(c *fiber.Ctx) error {
	var dto entity.{{.DomainPackage}}Model
	err := c.BodyParser(&dto)
	shared.Exception(err, fmt.Sprintf(message.MSG_PAYLOAD, entity.Name))

	errs := validate.ValidateStruct(dto)
	if len(errs) > 0 {
		return response.BadRequest(c, message.DATA_VALIDATION, errs)
	}

	ctx, cancel := shared.CreateContextWithTimeoutAndValue(c)
	defer cancel()

	userLog := c.Locals(constant.REQ_USERNAME).(string)
	err = h.feat.Delete(ctx, &dto.Code, userLog)
	shared.Exception(err, entity.Name)

	return response.OK(c, fmt.Sprintf(message.MSG_DELETED, entity.Name), nil)
}
`
