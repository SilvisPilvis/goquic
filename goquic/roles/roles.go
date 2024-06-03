package roles

type Role struct {
	Name        string
	Permissions []int
}

const (
	READ = iota
	WRITE
	EXECUTE
	SAVE
)

func (r *Role) Allow() string {

}
