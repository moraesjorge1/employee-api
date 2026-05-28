package employee

import (
	"errors"
	"testing"
)

type mockRepo struct {
	createFn      func(Employee) (Employee, error)
	getAllFn       func() ([]Employee, error)
	getFilteredFn func(ReportFilter) ([]Employee, error)
	getByIDFn     func(int64) (Employee, error)
	updateFn      func(Employee) (Employee, error)
	deleteFn      func(int64) error
}

func (m *mockRepo) Create(emp Employee) (Employee, error)          { return m.createFn(emp) }
func (m *mockRepo) GetAll() ([]Employee, error)                    { return m.getAllFn() }
func (m *mockRepo) GetFiltered(f ReportFilter) ([]Employee, error) { return m.getFilteredFn(f) }
func (m *mockRepo) GetByID(id int64) (Employee, error)             { return m.getByIDFn(id) }
func (m *mockRepo) Update(emp Employee) (Employee, error)          { return m.updateFn(emp) }
func (m *mockRepo) Delete(id int64) error                          { return m.deleteFn(id) }

func newService(repo EmployeeRepository) *Service {
	return NewService(repo)
}

// --- Create ---

func TestCreate_Valid(t *testing.T) {
	repo := &mockRepo{
		createFn: func(emp Employee) (Employee, error) {
			emp.ID = 1
			return emp, nil
		},
	}
	svc := newService(repo)

	emp, err := svc.Create(Employee{Name: "Ana", Position: "Dev", Salary: 5000, Type: "fulltime"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if emp.ID != 1 {
		t.Errorf("expected ID 1, got %d", emp.ID)
	}
}

func TestCreate_EmptyName(t *testing.T) {
	svc := newService(&mockRepo{})
	_, err := svc.Create(Employee{Position: "Dev", Salary: 5000, Type: "fulltime"})
	if err == nil || err.Error() != "name cannot be empty" {
		t.Errorf("expected 'name cannot be empty', got %v", err)
	}
}

func TestCreate_EmptyPosition(t *testing.T) {
	svc := newService(&mockRepo{})
	_, err := svc.Create(Employee{Name: "Ana", Salary: 5000, Type: "fulltime"})
	if err == nil || err.Error() != "position cannot be empty" {
		t.Errorf("expected 'position cannot be empty', got %v", err)
	}
}

func TestCreate_InvalidSalary(t *testing.T) {
	svc := newService(&mockRepo{})
	_, err := svc.Create(Employee{Name: "Ana", Position: "Dev", Salary: -1, Type: "fulltime"})
	if err == nil || err.Error() != "salary must be greater than zero" {
		t.Errorf("expected salary error, got %v", err)
	}
}

func TestCreate_InvalidType(t *testing.T) {
	svc := newService(&mockRepo{})
	_, err := svc.Create(Employee{Name: "Ana", Position: "Dev", Salary: 5000, Type: "intern"})
	if err == nil || err.Error() != "type must be 'fulltime' or 'contractor'" {
		t.Errorf("expected type error, got %v", err)
	}
}

// --- GetByID ---

func TestGetByID_Valid(t *testing.T) {
	want := Employee{ID: 1, Name: "Ana", Position: "Dev", Salary: 5000, Type: "fulltime"}
	repo := &mockRepo{getByIDFn: func(id int64) (Employee, error) { return want, nil }}
	svc := newService(repo)

	got, err := svc.GetByID(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestGetByID_ZeroID(t *testing.T) {
	svc := newService(&mockRepo{})
	_, err := svc.GetByID(0)
	if err == nil || err.Error() != "id cannot be zero" {
		t.Errorf("expected 'id cannot be zero', got %v", err)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	repo := &mockRepo{getByIDFn: func(id int64) (Employee, error) { return Employee{}, ErrNotFound }}
	svc := newService(repo)

	_, err := svc.GetByID(99)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

// --- Update ---

func TestUpdate_Valid(t *testing.T) {
	emp := Employee{ID: 1, Name: "Ana", Position: "Dev", Salary: 6000, Type: "contractor"}
	repo := &mockRepo{updateFn: func(e Employee) (Employee, error) { return e, nil }}
	svc := newService(repo)

	got, err := svc.Update(emp)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.Salary != 6000 {
		t.Errorf("expected salary 6000, got %v", got.Salary)
	}
}

func TestUpdate_ZeroID(t *testing.T) {
	svc := newService(&mockRepo{})
	_, err := svc.Update(Employee{Name: "Ana", Position: "Dev", Salary: 5000, Type: "fulltime"})
	if err == nil || err.Error() != "id cannot be zero" {
		t.Errorf("expected 'id cannot be zero', got %v", err)
	}
}

func TestUpdate_InvalidType(t *testing.T) {
	svc := newService(&mockRepo{})
	_, err := svc.Update(Employee{ID: 1, Name: "Ana", Position: "Dev", Salary: 5000, Type: "freelancer"})
	if err == nil || err.Error() != "type must be 'fulltime' or 'contractor'" {
		t.Errorf("expected type error, got %v", err)
	}
}

// --- Delete ---

func TestDelete_Valid(t *testing.T) {
	repo := &mockRepo{deleteFn: func(id int64) error { return nil }}
	svc := newService(repo)

	if err := svc.Delete(1); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestDelete_ZeroID(t *testing.T) {
	svc := newService(&mockRepo{})
	if err := svc.Delete(0); err == nil || err.Error() != "id cannot be zero" {
		t.Errorf("expected 'id cannot be zero', got %v", err)
	}
}

func TestDelete_NotFound(t *testing.T) {
	repo := &mockRepo{deleteFn: func(id int64) error { return ErrNotFound }}
	svc := newService(repo)

	if err := svc.Delete(99); !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

// --- GenerateReport ---

func TestGenerateReport_Empty(t *testing.T) {
	repo := &mockRepo{
		getFilteredFn: func(ReportFilter) ([]Employee, error) { return []Employee{}, nil },
	}
	svc := newService(repo)

	report, err := svc.GenerateReport(ReportFilter{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if report.TotalEmployees != 0 || report.AverageSalary != 0 {
		t.Errorf("expected empty report, got %+v", report)
	}
}

func TestGenerateReport_WithEmployees(t *testing.T) {
	employees := []Employee{
		{ID: 1, Name: "Ana", Position: "Dev", Salary: 4000, Type: "fulltime"},
		{ID: 2, Name: "Bob", Position: "QA", Salary: 6000, Type: "contractor"},
		{ID: 3, Name: "Lia", Position: "Dev", Salary: 2000, Type: "fulltime"},
	}
	repo := &mockRepo{
		getFilteredFn: func(ReportFilter) ([]Employee, error) { return employees, nil },
	}
	svc := newService(repo)

	report, err := svc.GenerateReport(ReportFilter{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if report.TotalEmployees != 3 {
		t.Errorf("expected 3 employees, got %d", report.TotalEmployees)
	}
	if report.TotalSalary != 12000 {
		t.Errorf("expected total salary 12000, got %v", report.TotalSalary)
	}
	if report.AverageSalary != 4000 {
		t.Errorf("expected average 4000, got %v", report.AverageSalary)
	}
	if report.ByType["fulltime"] != 2 || report.ByType["contractor"] != 1 {
		t.Errorf("unexpected ByType: %v", report.ByType)
	}
}

func TestGenerateReport_FilterByType(t *testing.T) {
	var capturedFilter ReportFilter
	repo := &mockRepo{
		getFilteredFn: func(f ReportFilter) ([]Employee, error) {
			capturedFilter = f
			return []Employee{
				{ID: 1, Name: "Ana", Position: "Dev", Salary: 4000, Type: "fulltime"},
			}, nil
		},
	}
	svc := newService(repo)

	report, err := svc.GenerateReport(ReportFilter{Type: "fulltime"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if capturedFilter.Type != "fulltime" {
		t.Errorf("expected filter type 'fulltime', got '%s'", capturedFilter.Type)
	}
	if report.TotalEmployees != 1 {
		t.Errorf("expected 1 employee, got %d", report.TotalEmployees)
	}
}

func TestGenerateReport_FilterBySalaryRange(t *testing.T) {
	var capturedFilter ReportFilter
	repo := &mockRepo{
		getFilteredFn: func(f ReportFilter) ([]Employee, error) {
			capturedFilter = f
			return []Employee{
				{ID: 2, Name: "Bob", Position: "QA", Salary: 6000, Type: "contractor"},
			}, nil
		},
	}
	svc := newService(repo)

	report, err := svc.GenerateReport(ReportFilter{MinSalary: 5000, MaxSalary: 8000})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if capturedFilter.MinSalary != 5000 || capturedFilter.MaxSalary != 8000 {
		t.Errorf("unexpected filter: %+v", capturedFilter)
	}
	if report.TotalSalary != 6000 {
		t.Errorf("expected total salary 6000, got %v", report.TotalSalary)
	}
}
