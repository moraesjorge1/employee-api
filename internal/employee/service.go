package employee

import "errors"

type EmployeeRepository interface {
	Create(emp Employee) (Employee, error)
	GetAll() ([]Employee, error)
	GetByID(id int64) (Employee, error)
	Update(emp Employee) (Employee, error)
	Delete(id int64) error
}

type Service struct {
	repo EmployeeRepository
}

func NewService(repo EmployeeRepository) *Service {
	return &Service{repo: repo}
}

var validTypes = map[string]bool{"fulltime": true, "contractor": true}

func (s *Service) Create(emp Employee) (Employee, error) {
	if emp.Name == "" {
		return Employee{}, errors.New("name cannot be empty")
	}
	if emp.Position == "" {
		return Employee{}, errors.New("position cannot be empty")
	}
	if emp.Salary <= 0 {
		return Employee{}, errors.New("salary must be greater than zero")
	}
	if !validTypes[emp.Type] {
		return Employee{}, errors.New("type must be 'fulltime' or 'contractor'")
	}
	return s.repo.Create(emp)
}

func (s *Service) GetAll() ([]Employee, error) {
	employees, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return employees, nil
}

func (s *Service) GetByID(id int64) (Employee, error) {
	if id == 0 {
		return Employee{}, errors.New("id cannot be zero")
	}
	return s.repo.GetByID(id)
}

func (s *Service) Update(emp Employee) (Employee, error) {
	if emp.ID == 0 {
		return Employee{}, errors.New("id cannot be zero")
	}
	if emp.Name == "" {
		return Employee{}, errors.New("name cannot be empty")
	}
	if emp.Position == "" {
		return Employee{}, errors.New("position cannot be empty")
	}
	if emp.Salary <= 0 {
		return Employee{}, errors.New("salary must be greater than zero")
	}
	if !validTypes[emp.Type] {
		return Employee{}, errors.New("type must be 'fulltime' or 'contractor'")
	}
	return s.repo.Update(emp)
}

func (s *Service) Delete(id int64) error {
	if id == 0 {
		return errors.New("id cannot be zero")
	}
	return s.repo.Delete(id)
}

func (s *Service) GenerateReport() (Reporting, error) {

	employees, err := s.repo.GetAll()
	if err != nil {
		return Reporting{}, err
	}

	report := Reporting{
		ByType:    make(map[string]int),
		Employees: employees,
	}

	for _, emp := range employees {
		report.TotalSalary += emp.Salary
		report.TotalEmployees++
		report.ByType[emp.Type]++
	}

	if report.TotalEmployees > 0 {
		report.AverageSalary = report.TotalSalary / float64(report.TotalEmployees)
	}

	return report, nil
}
