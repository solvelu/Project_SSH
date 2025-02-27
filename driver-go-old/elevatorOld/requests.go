package elevator

import (
	"Driver-go/elevio"
	"fmt"
)


type DirnBehaviourPair struct{
	Direction elevio.MotorDirection
	Behaviour ElevatorBehaviour
}

func requestsAbove(e Elevator) bool {
	for f := e.Floor + 1; f < NumFloors; f++ {
		for b := 0; b < NumButtons; b++ {
			if e.Requests[f][b] {
				return true
			}
		}
	}
	return false
}

func requestsBelow(e Elevator) bool {
	for f := 0; f < e.Floor; f++ {
		for b := 0; b < NumButtons; b++ {
			if e.Requests[f][b] {
				return true
			}
		}
	}
	return false
}

func requestsHere(e Elevator) bool {
	for b := 0; b < NumButtons; b++ {
		if e.Requests[e.Floor][b] {
			return true
		}
	}
	return false
}

func requestsChooseDirection(e Elevator) DirnBehaviourPair{
	switch e.Dirn {	
	case elevio.MD_Up:	
		if requestsAbove(e) {				
			return DirnBehaviourPair{Direction: elevio.MD_Up, Behaviour: EB_moving}
		}
		if requestsBelow(e) {
			return DirnBehaviourPair{Direction: elevio.MD_Down, Behaviour:EB_moving}
		}
		if requestsHere(e) {
			return DirnBehaviourPair{Direction: elevio.MD_Stop, Behaviour: EB_doorOpen}
		}
		return DirnBehaviourPair{Direction: elevio.MD_Stop, Behaviour: EB_idle}
	case elevio.MD_Down:
		if requestsBelow(e) {
			return DirnBehaviourPair{Direction: elevio.MD_Down, Behaviour: EB_moving}
		}
		if requestsAbove(e) {
			return DirnBehaviourPair{Direction: elevio.MD_Up, Behaviour: EB_moving}
		}
		if requestsHere(e) {
			return DirnBehaviourPair{Direction: elevio.MD_Stop, Behaviour: EB_doorOpen}
		}
		return DirnBehaviourPair{Direction: elevio.MD_Stop, Behaviour: EB_idle}
	case elevio.MD_Stop:		
		if requestsAbove(e) {
			return DirnBehaviourPair{Direction: elevio.MD_Up, Behaviour: EB_moving}
		}
		if requestsBelow(e) {
			return DirnBehaviourPair{Direction: elevio.MD_Down, Behaviour: EB_moving}
		}	
		if requestsHere(e) {
			fmt.Printf("Requests here\n")
			return DirnBehaviourPair{Direction: elevio.MD_Stop, Behaviour: EB_doorOpen}
		}
		return DirnBehaviourPair{Direction: elevio.MD_Stop, Behaviour: EB_idle}
		
	default:
		return DirnBehaviourPair{Direction: elevio.MD_Stop, Behaviour: EB_idle}
	}
		
}
func requestsShouldStop(e Elevator) bool {
	switch e.Dirn {
	case elevio.MD_Up:
		return e.Requests[e.Floor][elevio.BT_HallUp] || e.Requests[e.Floor][elevio.BT_Cab] || !requestsAbove(e)
	case elevio.MD_Down:
		return e.Requests[e.Floor][elevio.BT_HallDown] || e.Requests[e.Floor][elevio.BT_Cab] || !requestsBelow(e)
	case elevio.MD_Stop:
		return true
	default:
		return false
	}
}
/*
func requestsShouldClearImmediately(e Elevator, btn_floor int, btn_type elevio.ButtonType) bool {
	switch e.Config.clearRequestVariant {
	case CRV_all:
		return e.Floor == btn_floor
	case CRV_InDirn:
		return e.Floor == btn_floor && 
		(
			(e.Dirn == elevio.MD_Up && btn_type == elevio.BT_HallUp) ||
		(e.Dirn == elevio.MD_Down && btn_type == elevio.BT_HallDown) ||
		e.Dirn == elevio.MD_Stop || btn_type == elevio.BT_Cab)
	default:
		return false
	}
}
*/

func requestsShouldClearImmediately(e Elevator, btn_floor int, btn_type elevio.ButtonType) bool {
    fmt.Printf("Checking if request should clear immediately: Floor=%d, BtnFloor=%d, BtnType=%d, Dirn=%d, ClearRequestVariant=%d\n", e.Floor, btn_floor, btn_type, e.Dirn, e.Config.clearRequestVariant)
    
    switch e.Config.clearRequestVariant {
    case CRV_all:
        result := (e.Floor == btn_floor)
        fmt.Printf("CRV_All: Floor=%d, BtnFloor=%d, Result=%d\n", e.Floor, btn_floor, result)
        return result
    case CRV_InDirn:
        result := (e.Floor == btn_floor && 
        (
            (e.Dirn == elevio.MD_Up && btn_type == elevio.BT_HallUp) ||
            (e.Dirn == elevio.MD_Down && btn_type == elevio.BT_HallDown) ||
            e.Dirn == elevio.MD_Stop || btn_type == elevio.BT_Cab))
			fmt.Printf("CRV_InDirn: Floor=%d, BtnFloor=%d, Dirn=%d, BtnType=%d, Result=%v\n", e.Floor, btn_floor, e.Dirn, btn_type, result)
        return result
    default:
        fmt.Printf("Default case: false\n")
        return false
    }
}

func ClearAtCurrentFloor(e Elevator) Elevator{
	switch e.Config.clearRequestVariant{
		case CRV_all:
			for b := 0; b < NumButtons; b++ {	
				e.Requests[e.Floor][b] = false
			}
		case CRV_InDirn:
			e.Requests[e.Floor][elevio.BT_Cab] = false
			switch e.Dirn {
			case elevio.MD_Up:
				if !requestsAbove(e) && !e.Requests[e.Floor][elevio.BT_HallUp] {
					e.Requests[e.Floor][elevio.BT_HallDown] = false
				}
				e.Requests[e.Floor][elevio.BT_HallUp] = false
	
			case elevio.MD_Down:
				if !requestsBelow(e) && !e.Requests[e.Floor][elevio.BT_HallDown] {
					e.Requests[e.Floor][elevio.BT_HallUp] = false
				}
				e.Requests[e.Floor][elevio.BT_HallDown] = false
	
			case elevio.MD_Stop:
			default:
				e.Requests[e.Floor][elevio.BT_HallUp] = false
				e.Requests[e.Floor][elevio.BT_HallDown] = false
			}
	}
	return e
}

