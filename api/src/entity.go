type User struct {
    Username      string
    Email         string
    BirthDate     Time
    JoinDate      Time
    DisplayPic    Image
    backgroundPic Image
    Followers     []FollowerView
}

type FollowerView struct {
    Username       string
    MiniDisplayPic Image
    Circle         Circle
}

type Circle int

const (
    Gold     Circle = iota
    Standard Circle = iota
)

