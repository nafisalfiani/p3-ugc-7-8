package handler

import (
	"api-gateway/entity"
	"api-gateway/errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ListUsers returns list of user
//
// @Summary Get user list
// @Description Returns list of user
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} entity.HttpResp{data=[]entity.User}
// @Failure 400 {object} entity.HttpResp
// @Failure 500 {object} entity.HttpResp
// @Router /v1/users [get]
func (h *Handler) ListUsers(c echo.Context) error {
	users, err := h.user.List(c.Request().Context())
	if err != nil {
		return h.httpError(c, err)
	}

	return h.httpSuccess(c, http.StatusOK, users)
}

// CreateUser creates new user
//
// @Summary Create user
// @Description Creates new user
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user body entity.UserCreateRequest true "user create request"
// @Success 200 {object} entity.HttpResp{data=entity.User}
// @Failure 400 {object} entity.HttpResp
// @Failure 500 {object} entity.HttpResp
// @Router /v1/users [post]
func (h *Handler) CreateUser(c echo.Context) error {
	req := entity.UserCreateRequest{}
	if err := c.Bind(&req); err != nil {
		return h.httpError(c, err)
	}

	if err := h.validator.Struct(req); err != nil {
		return h.httpError(c, errors.ErrBadRequest, err.Error())
	}

	user := entity.User{
		Name:  req.Name,
		Email: req.Email,
	}
	newUser, err := h.user.Create(c.Request().Context(), user)
	if err != nil {
		return h.httpError(c, err)
	}

	return h.httpSuccess(c, http.StatusOK, newUser)
}

// GetUsers returns specific user
//
// @Summary Get user detail
// @Description Get specific user
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "user id"
// @Success 200 {object} entity.HttpResp{data=entity.User}
// @Failure 400 {object} entity.HttpResp
// @Failure 500 {object} entity.HttpResp
// @Router /v1/users/{id} [get]
func (h *Handler) GetUser(c echo.Context) error {
	req := entity.UserGetRequest{}
	if err := c.Bind(&req); err != nil {
		return h.httpError(c, err)
	}

	if err := h.validator.Struct(req); err != nil {
		return h.httpError(c, errors.ErrBadRequest, err.Error())
	}

	user, err := h.user.Get(c.Request().Context(), entity.User{Id: req.Id})
	if err != nil {
		return h.httpError(c, err)
	}

	return h.httpSuccess(c, http.StatusOK, user)
}

// UpdateUser returns logged in user detail
//
// @Summary Get user detail
// @Description Get logged in user detail
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "user id"
// @Param user body entity.UserUpdateRequest true "user update request"
// @Success 200 {object} entity.HttpResp{data=entity.User}
// @Failure 400 {object} entity.HttpResp
// @Failure 500 {object} entity.HttpResp
// @Router /v1/users/{id} [put]
func (h *Handler) UpdateUser(c echo.Context) error {
	req := entity.UserUpdateRequest{}
	if err := c.Bind(&req); err != nil {
		return h.httpError(c, err)
	}

	if err := h.validator.Struct(req); err != nil {
		return h.httpError(c, errors.ErrBadRequest, err.Error())
	}

	user := entity.User{
		Id:    req.Id,
		Name:  req.Name,
		Email: req.Email,
	}
	user, err := h.user.Update(c.Request().Context(), user)
	if err != nil {
		return h.httpError(c, err)
	}

	return h.httpSuccess(c, http.StatusOK, user)
}

// DeleteUser deletes existing user
//
// @Summary Delete user
// @Description delete existing user
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "user id"
// @Success 200 {object} entity.HttpResp
// @Failure 400 {object} entity.HttpResp
// @Failure 500 {object} entity.HttpResp
// @Router /v1/users/{id} [delete]
func (h *Handler) DeleteUser(c echo.Context) error {
	req := entity.UserGetRequest{}
	if err := c.Bind(&req); err != nil {
		return h.httpError(c, err)
	}

	if err := h.validator.Struct(req); err != nil {
		return h.httpError(c, errors.ErrBadRequest, err.Error())
	}

	if err := h.user.Delete(c.Request().Context(), entity.User{Id: req.Id}); err != nil {
		return h.httpError(c, err)
	}

	return h.httpSuccess(c, http.StatusOK, nil)
}
