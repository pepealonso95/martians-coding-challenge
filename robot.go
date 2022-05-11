package main

import (
	"errors"
	"strconv"
)

const orientationLimit = 4

const (
	Left  int = -1
	Right     = 1
)

const (
	North int = 0
	East      = 1
	South     = 2
	West      = 3
)

// Representation of a robot
type Robot struct {
	// Whether robot is lost or not
	Lost bool
	// Where the robot is pointing to (N,E,S,W)
	Orientation int
	// Latitude coordinate representation of robot position
	X int
	// Longitude coordinate representation of robot position
	Y int
}

// Returns a string representation of the robot's position
func (robot *Robot) GetPos() string {
	return strconv.Itoa(robot.X) + " " + strconv.Itoa(robot.Y)
}

// Decides how to process the given instruction for the robot
func (robot *Robot) Instruct(direction string, env *Environment) (bool, error) {
	switch direction {
	case "F":
		robot.Move(env)
		return false, nil
	case "L":
		robot.Turn(Left)
		return false, nil
	case "R":
		robot.Turn(Right)
		return false, nil
	default:
		return true, errors.New("Invalid instruction")
	}
}

// Function to turn right
func (robot *Robot) Turn(move int) {
	robot.Orientation = (robot.Orientation + move + orientationLimit) % orientationLimit
}

// Function set robot orientation from string
func (robot *Robot) setOrientation(orientation string) error {
	switch orientation {
	case "N":
		robot.Orientation = North
		return nil
	case "E":
		robot.Orientation = East
		return nil
	case "W":
		robot.Orientation = West
		return nil
	case "S":
		robot.Orientation = South
		return nil
	default:
		return errors.New("Invalid orientation")
	}
}

// Function get string representation of orientation
func (robot *Robot) getOrientation() (string, error) {
	switch robot.Orientation {
	case North:
		return "N", nil
	case East:
		return "E", nil
	case West:
		return "W", nil
	case South:
		return "S", nil
	default:
		return "", errors.New("Invalid orientation")
	}
}

// Checks in which direction and lets the robot advance if no lost robot is on its way
func (robot *Robot) Move(env *Environment) (bool, error) {
	x := robot.X
	y := robot.Y
	switch robot.Orientation {
	case North:
		y++
		break
	case East:
		x++
		break
	case West:
		x--
		break
	case South:
		y--
		break
	default:
		return true, errors.New("Invalid orientation")
	}
	if env.checkLostRobot(x, y) {
		return true, nil
	}
	return robot.advance(x, y, env), nil
}

// Prints last known position, orientation of robot and whether its is lost or not
func (robot *Robot) printRobotPos() string {
	robotLostString := ""
	if robot.Lost {
		robotLostString = "LOST"
	}
	robotOrienation, err := robot.getOrientation()
	if err != nil {
		return "Error"
	}
	return strconv.Itoa(robot.X) + " " + strconv.Itoa(robot.Y) + " " + robotOrienation + " " + robotLostString
}

// Advance the robot to the next position or report it as lost and keep its last seen position
func (robot *Robot) advance(x int, y int, env *Environment) bool {
	if y < 0 || x < 0 || y > env.limitY || x > env.limitX {
		robot.Lost = true
		env.addLostRobot(x, y)
		return true
	} else {
		robot.X = x
		robot.Y = y
		return false
	}
}
