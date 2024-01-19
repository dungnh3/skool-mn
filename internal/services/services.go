package services

import (
	errorz "github.com/dungnh3/skool-mn/internal/errors"
	"github.com/dungnh3/skool-mn/internal/models"
	store "github.com/dungnh3/skool-mn/internal/models/store"
	"github.com/dungnh3/skool-mn/internal/repositories"
	"github.com/dungnh3/skool-mn/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// @Summary list all students of parent
// @Description list all students of parent
// @Tags Parents
// @Accept json
// @Produce json
// @Param id path string true "parent_id"
// @Success 200 {object} models.Account
// @Router /parents/{id}/students [get]
func (s *Server) listStudents(ctx *gin.Context) {
	timeoutCtx, cancel := createTimeoutContext(ctx)
	defer cancel()

	parentId := getIDFromPath(ctx)
	accounts, err := s.r.ListStudents(timeoutCtx, parentId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorz.NewErrResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}

type (
	RegisterPickUpTime struct {
		ParentId     string   `json:"parent_id" binding:"required"`
		StudentIds   []string `json:"student_ids" binding:"required"`
		RegisterTime string   `json:"register_time" binding:"required"`
	}
)

// @Summary register a new schedule pick up
// @Description register a new schedule pick up
// @Tags Registers
// @Accept json
// @Produce json
// @Param register body RegisterPickUpTime true "register a new schedule"
// @Success 200 {object} models.Register
// @Router /parents/register [post]
func (s *Server) registerPickUpTime(ctx *gin.Context) {
	timeoutCtx, cancel := createTimeoutContext(ctx)
	defer cancel()

	var request RegisterPickUpTime
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorz.NewErrResponse(errorz.ErrBadParamInput, err.Error()))
		return
	}

	registerTime, err := time.Parse("2006-01-02 15:04:05", request.RegisterTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorz.NewErrResponse(errorz.ErrTimeRegisterInputInvalid, err.Error()))
		return
	}

	var registers []*models.Register
	for idx := range request.StudentIds {
		registers = append(registers, &models.Register{
			Register: store.Register{
				ID:           utils.GenerateId(),
				ParentID:     request.ParentId,
				StudentID:    request.StudentIds[idx],
				Status:       store.Registered,
				RegisterTime: registerTime,
			},
		})
	}

	if err = s.r.Transaction(func(r repositories.Repository) error {
		for idx := range registers {
			if err = r.CreateRegister(timeoutCtx, registers[idx]); err != nil {
				return err
			}

			tx := models.Transaction{
				Transaction: store.Transaction{
					RegisterID: registers[idx].ID,
					ActionType: store.ActionRegis,
				},
			}
			if err = r.CreateTransaction(timeoutCtx, &tx); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		ctx.JSON(http.StatusBadRequest, errorz.NewErrResponse(errorz.ErrInternalServerError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, registers)
}

// @Summary confirm that parents is ready for pik up
// @Description confirm that parents is ready for pik up
// @Tags Registers
// @Accept json
// @Produce json
// @Param id path string true "register_id"
// @Success 200 {object} models.Register
// @Router /parents/registers/:id/waiting [put]
func (s *Server) waitingFromParent(ctx *gin.Context) {
	timeoutCtx, cancel := createTimeoutContext(ctx)
	defer cancel()

	registerId := getIDFromPath(ctx)
	if err := s.r.Transaction(func(r repositories.Repository) error {
		lockedRegister, err := r.LockRegister(timeoutCtx, registerId)
		if err != nil {
			return err
		}

		if err := r.WaitingFromParent(timeoutCtx, lockedRegister.ID); err != nil {
			return err
		}

		tx := models.Transaction{
			Transaction: store.Transaction{
				RegisterID: lockedRegister.ID,
				ActionType: store.ActionWait,
			},
		}
		if err := r.CreateTransaction(timeoutCtx, &tx); err != nil {
			return err
		}
		return nil
	}); err != nil {
		ctx.JSON(http.StatusBadRequest, errorz.NewErrResponse(errorz.ErrInternalServerError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, "OK")
}

// @Summary confirm that parents want to cancel register
// @Description confirm that parents want to cancel register
// @Tags Registers
// @Accept json
// @Produce json
// @Param id path string true "register_id"
// @Success 200 {object} models.Register
// @Router /parents/registers/:id/cancel [put]
func (s *Server) cancelFromParent(ctx *gin.Context) {
	timeoutCtx, cancel := createTimeoutContext(ctx)
	defer cancel()

	registerId := getIDFromPath(ctx)
	if err := s.r.Transaction(func(r repositories.Repository) error {
		lockedRegister, err := r.LockRegister(timeoutCtx, registerId)
		if err != nil {
			return err
		}

		if err := r.CancelFromParent(timeoutCtx, lockedRegister.ID); err != nil {
			return err
		}

		tx := models.Transaction{
			Transaction: store.Transaction{
				RegisterID: lockedRegister.ID,
				ActionType: store.ActionCancel,
			},
		}
		if err := r.CreateTransaction(timeoutCtx, &tx); err != nil {
			return err
		}
		return nil
	}); err != nil {
		ctx.JSON(http.StatusBadRequest, errorz.NewErrResponse(errorz.ErrInternalServerError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, "OK")
}

// @Summary confirm that parents have completed
// @Description confirm that parents have completed
// @Tags Registers
// @Accept json
// @Produce json
// @Param id path string true "register_id"
// @Success 200 {object} models.Register
// @Router /parents/registers/:id/confirm [put]
func (s *Server) confirmCompleted(ctx *gin.Context) {
	timeoutCtx, cancel := createTimeoutContext(ctx)
	defer cancel()

	registerId := getIDFromPath(ctx)
	if err := s.r.Transaction(func(r repositories.Repository) error {
		lockedRegister, err := r.LockRegister(timeoutCtx, registerId)
		if err != nil {
			return err
		}

		if err := r.ConfirmFromParent(timeoutCtx, lockedRegister.ID); err != nil {
			return err
		}

		tx := models.Transaction{
			Transaction: store.Transaction{
				RegisterID: lockedRegister.ID,
				ActionType: store.ActionParentConfirm,
			},
		}
		if err := r.CreateTransaction(timeoutCtx, &tx); err != nil {
			return err
		}
		return nil
	}); err != nil {
		ctx.JSON(http.StatusBadRequest, errorz.NewErrResponse(errorz.ErrInternalServerError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, "OK")
}

// @Summary confirm that their student has left the class
// @Description confirm that their student has left the class
// @Tags Teachers
// @Accept json
// @Produce json
// @Param id path string true "register_id"
// @Success 200 {object} models.Register
// @Router /teachers/registers/:id/confirm [put]
func (s *Server) confirmFromTeacher(ctx *gin.Context) {
	timeoutCtx, cancel := createTimeoutContext(ctx)
	defer cancel()

	registerId := getIDFromPath(ctx)
	if err := s.r.Transaction(func(r repositories.Repository) error {
		lockedRegister, err := r.LockRegister(timeoutCtx, registerId)
		if err != nil {
			return err
		}

		if err := r.ConfirmFromTeacher(timeoutCtx, lockedRegister.ID); err != nil {
			return err
		}

		tx := models.Transaction{
			Transaction: store.Transaction{
				RegisterID: lockedRegister.ID,
				ActionType: store.ActionConfirm,
			},
		}
		if err := r.CreateTransaction(timeoutCtx, &tx); err != nil {
			return err
		}
		return nil
	}); err != nil {
		ctx.JSON(http.StatusBadRequest, errorz.NewErrResponse(errorz.ErrInternalServerError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, "OK")
}

// @Summary confirm that this teacher want to reject the register schedule
// @Description confirm that this teacher want to reject the register schedule
// @Tags Teachers
// @Accept json
// @Produce json
// @Param id path string true "register_id"
// @Success 200 {object} models.Register
// @Router /teachers/registers/:id/reject [put]
func (s *Server) rejectFromTeacher(ctx *gin.Context) {
	timeoutCtx, cancel := createTimeoutContext(ctx)
	defer cancel()

	registerId := getIDFromPath(ctx)
	if err := s.r.Transaction(func(r repositories.Repository) error {
		lockedRegister, err := r.LockRegister(timeoutCtx, registerId)
		if err != nil {
			return err
		}

		if err := r.RejectFromTeacher(timeoutCtx, lockedRegister.ID); err != nil {
			return err
		}

		tx := models.Transaction{
			Transaction: store.Transaction{
				RegisterID: lockedRegister.ID,
				ActionType: store.ActionReject,
			},
		}
		if err := r.CreateTransaction(timeoutCtx, &tx); err != nil {
			return err
		}
		return nil
	}); err != nil {
		ctx.JSON(http.StatusBadRequest, errorz.NewErrResponse(errorz.ErrInternalServerError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, "OK")
}

// @Summary confirm that this student left this class
// @Description confirm that this teacher want to reject the register schedule
// @Tags Students
// @Accept json
// @Produce json
// @Param id path string true "student_id"
// @Success 200
// @Router /students/:id/leave [put]
func (s *Server) studentLeaveClass(ctx *gin.Context) {
	timeoutCtx, cancel := createTimeoutContext(ctx)
	defer cancel()

	studentId := getIDFromPath(ctx)
	register, err := s.r.GetConfirmedRegisterLatest(timeoutCtx, studentId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorz.NewErrResponse(errorz.ErrExceptionError, err.Error()))
		return
	}

	if err := s.r.Transaction(func(r repositories.Repository) error {
		lockedRegister, err := r.LockRegister(timeoutCtx, register.ID)
		if err != nil {
			return err
		}

		if err := r.StudentLeaveClass(timeoutCtx, lockedRegister.ID); err != nil {
			return err
		}

		tx := models.Transaction{
			Transaction: store.Transaction{
				RegisterID: lockedRegister.ID,
				ActionType: store.ActionLeaveClass,
			},
		}
		if err := r.CreateTransaction(timeoutCtx, &tx); err != nil {
			return err
		}
		return nil
	}); err != nil {
		ctx.JSON(http.StatusBadRequest, errorz.NewErrResponse(errorz.ErrInternalServerError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, "OK")
}

// @Summary confirm that this student left this school
// @Description confirm that this student left this school
// @Tags Students
// @Accept json
// @Produce json
// @Param id path string true "student_id"
// @Success 200
// @Router /students/:id/out [put]
func (s *Server) studentOutSchool(ctx *gin.Context) {
	timeoutCtx, cancel := createTimeoutContext(ctx)
	defer cancel()

	studentId := getIDFromPath(ctx)
	register, err := s.r.GetLeftRegisterLatest(timeoutCtx, studentId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorz.NewErrResponse(errorz.ErrExceptionError, err.Error()))
		return
	}

	if err := s.r.Transaction(func(r repositories.Repository) error {
		lockedRegister, err := r.LockRegister(timeoutCtx, register.ID)
		if err != nil {
			return err
		}

		if err := r.StudentOutSchool(timeoutCtx, lockedRegister.ID); err != nil {
			return err
		}

		tx := models.Transaction{
			Transaction: store.Transaction{
				RegisterID: lockedRegister.ID,
				ActionType: store.ActionOutSchool,
			},
		}
		if err := r.CreateTransaction(timeoutCtx, &tx); err != nil {
			return err
		}
		return nil
	}); err != nil {
		ctx.JSON(http.StatusBadRequest, errorz.NewErrResponse(errorz.ErrInternalServerError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, "OK")
}
