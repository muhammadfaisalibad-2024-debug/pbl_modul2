package service

import (
	"errors"
	"strconv"

	"api-students/app/model"
	"api-students/app/repository"
	"api-students/helper"

	"github.com/gofiber/fiber/v2"
)

type SupervisorService struct {
	supervisors repository.SupervisorRepository
	tx          *repository.TxManager
}

func NewSupervisorService(supervisors repository.SupervisorRepository, tx *repository.TxManager) *SupervisorService {
	return &SupervisorService{supervisors: supervisors, tx: tx}
}

// POST /api/v1/supervisors/:id/assign
func (s *SupervisorService) Assign(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	supervisorID, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("id pembimbing harus berupa angka positif")
	}

	var req model.AssignSupervisorRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}
	if errs := helper.ValidateStruct(req); len(errs) > 0 {
		return helper.Validation(errs)
	}

	var updatedSupervisor model.Supervisor
	err := s.tx.WithTx(ctx, func(tx repository.DBTX) error {
		repo := s.supervisors.WithTx(tx)

		// 1. Kunci baris mahasiswa untuk mencegah penetapan serentak ganda
		if err := repo.LockStudent(ctx, req.StudentID); err != nil {
			return err
		}

		// 2. Increment kuota secara atomik
		sup, err := repo.IncrementSupervised(ctx, supervisorID)
		if err != nil {
			return err
		}
		updatedSupervisor = sup

		// 3. Catat bimbingan (UNIQUE student_id menjaga di level DB)
		return repo.InsertSupervision(ctx, req.StudentID, supervisorID)
	})

	if err != nil {
		return translateSupervisorError(err)
	}

	return helper.Created(c, "berhasil menetapkan dosen pembimbing", updatedSupervisor,
		"/api/v1/supervisors/"+strconv.Itoa(supervisorID))
}

// POST /api/v1/supervisions/switch
func (s *SupervisorService) Switch(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	var req model.SwitchSupervisorRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}
	if errs := helper.ValidateStruct(req); len(errs) > 0 {
		return helper.Validation(errs)
	}
	if req.OldSupervisorID == req.NewSupervisorID {
		return helper.Validation(map[string]string{
			"new_supervisor_id": "tidak boleh sama dengan old_supervisor_id",
		})
	}

	var newSupervisor model.Supervisor
	err := s.tx.WithTx(ctx, func(tx repository.DBTX) error {
		repo := s.supervisors.WithTx(tx)

		// 1. Kunci baris mahasiswa
		if err := repo.LockStudent(ctx, req.StudentID); err != nil {
			return err
		}

		// 2. Kunci KEDUA pembimbing dalam urutan id menaik (bebas deadlock)
		if err := repo.LockSupervisors(ctx, req.OldSupervisorID, req.NewSupervisorID); err != nil {
			return err
		}

		// 3. Hapus bimbingan lama
		if err := repo.DeleteSupervision(ctx, req.StudentID, req.OldSupervisorID); err != nil {
			return err
		}

		// 4. Kurangi kuota pembimbing lama
		if _, err := repo.DecrementSupervised(ctx, req.OldSupervisorID); err != nil {
			return err
		}

		// 5. Tambah kuota pembimbing baru secara atomik
		sup, err := repo.IncrementSupervised(ctx, req.NewSupervisorID)
		if err != nil {
			return err
		}
		newSupervisor = sup

		// 6. Catat bimbingan baru
		return repo.InsertSupervision(ctx, req.StudentID, req.NewSupervisorID)
	})

	if err != nil {
		return translateSupervisorError(err)
	}

	return helper.Success(c, fiber.StatusOK, "berhasil mengganti dosen pembimbing", newSupervisor)
}

// PATCH /api/v1/supervisors/:id
func (s *SupervisorService) UpdateQuota(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("id harus berupa angka positif")
	}

	var req model.UpdateSupervisorQuotaRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}
	if errs := helper.ValidateStruct(req); len(errs) > 0 {
		return helper.Validation(errs)
	}

	updated, err := s.supervisors.UpdateQuota(ctx, id, req.Quota, req.Version)
	if err != nil {
		return translateSupervisorError(err)
	}

	return helper.Success(c, fiber.StatusOK, "kuota pembimbing berhasil diubah", updated)
}

func translateSupervisorError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return helper.NotFound("data pembimbing atau mahasiswa tidak ditemukan")
	case errors.Is(err, repository.ErrSupervisorFull):
		return helper.NewError(fiber.StatusConflict, "SUPERVISOR_FULL", "kuota dosen pembimbing sudah penuh")
	case errors.Is(err, repository.ErrAlreadySupervised):
		return helper.NewError(fiber.StatusConflict, "ALREADY_SUPERVISED", "mahasiswa sudah memiliki pembimbing")
	case errors.Is(err, repository.ErrSupervisionNotFound):
		return helper.NewError(fiber.StatusConflict, "NOT_SUPERVISED", "mahasiswa tidak dibimbing oleh dosen tersebut")
	case errors.Is(err, repository.ErrVersionConflict):
		return helper.NewError(fiber.StatusConflict, "VERSION_CONFLICT",
			"data pembimbing sudah diubah pihak lain, muat ulang lalu coba lagi")
	case errors.Is(err, repository.ErrQuotaBelowFilled):
		return helper.NewError(fiber.StatusConflict, "QUOTA_BELOW_ENROLLED",
			"kuota tidak boleh lebih kecil dari jumlah mahasiswa yang dibimbing")
	default:
		return helper.Internal(err)
	}
}
