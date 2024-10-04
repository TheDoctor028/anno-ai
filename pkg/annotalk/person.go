package annotalk

type PersonGender string

const (
	Man      PersonGender = "man"
	Woman                 = "woman"
	Whatever              = "whatever"
)

type Person struct {
	Name               string       `json:"name"`
	Age                int          `json:"age"`
	Gender             PersonGender `json:"gender"`
	InterestedInGender PersonGender `json:"interestedInGender"`
	Description        string       `json:"description"`
}
