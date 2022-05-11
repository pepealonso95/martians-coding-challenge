package main

import (
	"errors"
	"strconv"
)

// Representation of a robot
type Robot struct {
	// Whether robot is lost or not
	Lost bool
	// Where the robot is pointing to (N,E,S,W)
	Orientation string
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
		return robot.Move(env)
	case "L":
		return robot.TurnLeft()
	case "R":
		return robot.TurnRight()
	default:
		return false, errors.New("Invalid instruction")
	}
}

// NOTE: There are smarter ways to do turn left and right, but the simplest possible (and most readable)
// are these two switch statement methods (two separate methods as code is not really being that reused here)

// In GO there is no generic library function to get index of value in slice, so
// getting the position with the indes +1 or -1 depending on the direction is not as simple as in say JS
// In JS it could be done using an ordered array of orientations [N,E,S,W] and doing:
// orientations[(orientations.indexOf(robot.orientation) + 1 (RIGHT) or -1 (LEFT) + 4) % 4]

// It doesnt matter that much though as the orientations are only 4 so the code looks clean enough

// Function to turn right
func (robot *Robot) TurnRight() (bool, error) {
	switch robot.Orientation {
	case "N":
		robot.Orientation = "E"
		return false, nil
	case "E":
		robot.Orientation = "S"
		return false, nil
	case "W":
		robot.Orientation = "N"
		return false, nil
	case "S":
		robot.Orientation = "W"
		return false, nil
	default:
		return true, errors.New("Invalid orientation")
	}
}

// Function to turn left
func (robot *Robot) TurnLeft() (bool, error) {
	switch robot.Orientation {
	case "N":
		robot.Orientation = "W"
		return false, nil
	case "E":
		robot.Orientation = "N"
		return false, nil
	case "W":
		robot.Orientation = "S"
		return false, nil
	case "S":
		robot.Orientation = "E"
		return false, nil
	default:
		return true, errors.New("Invalid orientation")
	}
}

// Checks in which direction and lets the robot advance if no lost robot is on its way
func (robot *Robot) Move(env *Environment) (bool, error) {
	switch robot.Orientation {
	case "N":
		if env.checkLostRobot(robot.X, robot.Y+1) {
			return false, nil
		}
		return robot.advance(robot.X, robot.Y+1, env), nil
	case "E":
		if env.checkLostRobot(robot.X+1, robot.Y) {
			return false, nil
		}
		return robot.advance(robot.X+1, robot.Y, env), nil
	case "W":
		if env.checkLostRobot(robot.X-1, robot.Y) {
			return false, nil
		}
		return robot.advance(robot.X-1, robot.Y, env), nil
	case "S":
		if env.checkLostRobot(robot.X, robot.Y-1) {
			return false, nil
		}
		return robot.advance(robot.X, robot.Y-1, env), nil
	default:
		return true, errors.New("Invalid orientation")
	}
}

// Prints last known position, orientation of robot and whether its is lost or not
func (robot *Robot) printRobotPos() string {
	robotLostString := ""
	if robot.Lost {
		robotLostString = "LOST"
	}
	return strconv.Itoa(robot.X) + " " + strconv.Itoa(robot.Y) + " " + robot.Orientation + " " + robotLostString
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
