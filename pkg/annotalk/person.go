package annotalk

type PersonGender string

const (
	Man      PersonGender = "man"
	Woman                 = "woman"
	Whatever              = "whatever"
)

type Person struct {
	Name               string
	Age                int
	Gender             PersonGender
	InterestedInGender PersonGender

	Description string
}
