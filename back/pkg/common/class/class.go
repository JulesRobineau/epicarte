package class

import (
	"gin-template/pkg/dto"
	"gin-template/pkg/model"
	error2 "gin-template/utils/error"
	"gorm.io/gorm"
)

// GetClassByID gets a class by ID
func GetClassByID(tx *gorm.DB, classId uint64) (*dto.Class, error) {
	classModel := model.NewClassModel(tx)
	class, err := classModel.GetByID(classId)
	if err != nil {
		return nil, error2.FromDatabaseError(err)
	}

	students := make([]dto.Student, 0)
	for _, student := range class.Students {
		students = append(students, dto.Student{
			ID:        student.ID,
			Email:     student.Email,
			FirstName: student.FirstName,
			LastName:  student.LastName,
		})
	}

	return &dto.Class{
		TinyClass: dto.TinyClass{
			ID:   class.ID,
			Name: class.Name,
			Year: class.Year,
		},
		Students: students,
	}, nil
}

// GetAllClasses gets all classes
func GetAllClasses(tx *gorm.DB) (*dto.ClassList, error) {
	classModel := model.NewClassModel(tx)
	classes, err := classModel.FindAll()
	if err != nil {
		return nil, error2.FromDatabaseError(err)
	}

	classDtos := make([]dto.TinyClass, 0)
	for _, class := range classes {
		classDtos = append(classDtos, dto.TinyClass{
			ID:   class.ID,
			Name: class.Name,
			Year: class.Year,
		})
	}

	return &dto.ClassList{Classes: classDtos}, nil
}

// CreateClass creates a new class
func CreateClass(tx *gorm.DB, class dto.CreateClass) (*dto.TinyClass, error) {
	classModel := model.NewClassModel(tx)
	cl := model.Class{
		Name: class.Name,
		Year: class.Year,
	}
	err := classModel.Create(&cl)
	if err != nil {
		return nil, error2.FromDatabaseError(err)
	}

	return &dto.TinyClass{
		ID:   cl.ID,
		Name: cl.Name,
		Year: cl.Year,
	}, nil
}

// UpdateClass updates a class
func UpdateClass(tx *gorm.DB, class dto.UpdateClass) error {
	classModel := model.NewClassModel(tx)
	cl := model.Class{
		ID:   class.Id,
		Name: class.Name,
		Year: class.Year,
	}
	err := classModel.Update(&cl)
	if err != nil {
		return error2.FromDatabaseError(err)
	}

	return nil
}

// DeleteClass deletes a class
func DeleteClass(tx *gorm.DB, classId uint64) error {
	classModel := model.NewClassModel(tx)
	err := classModel.Delete(classId)
	if err != nil {
		return error2.FromDatabaseError(err)
	}

	return nil
}

// AddStudentToClass adds a student to a class
func AddStudentToClass(tx *gorm.DB, student dto.AddStudentToClass) (*dto.Student, error) {
	studentModel := model.NewClassModel(tx)
	st := model.Student{
		Email:     student.Student.Email,
		FirstName: student.Student.FirstName,
		LastName:  student.Student.LastName,
		ClassID:   student.ClassId,
	}
	err := studentModel.AddStudent(&model.Class{ID: student.ClassId}, &st)
	if err != nil {
		return nil, error2.FromDatabaseError(err)
	}

	return &dto.Student{
		ID:        st.ID,
		Email:     st.Email,
		FirstName: st.FirstName,
		LastName:  st.LastName,
	}, nil
}
