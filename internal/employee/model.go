package employee

type Employee struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Position string  `json:"position"`
	Salary   float64 `json:"salary"`
	Type     string  `json:"type"` // "fulltime" ou "contractor"
}

type Reporting struct {
	TotalEmployees int            `json:"total_employees"`
	TotalSalary    float64        `json:"total_salary"`
	AverageSalary  float64        `json:"average_salary"`
	ByType         map[string]int `json:"by_type"`
	Employees      []Employee     `json:"employees"`
}
