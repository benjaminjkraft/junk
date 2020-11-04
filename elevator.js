{
    init: function(elevators, floors) {
        const floorStates = floors.map(() => false);
        elevators[0].goToFloor(10)
        for (let elevator of elevators) {
        let dir = "up";
        elevator.on("idle", function() {
            const cur = elevator.currentFloor()
            if (elevator === elevators[2]) {
                console.log(elevator.destinationQueue, cur, dir, floorStates);
            }
            if (dir === "up") {
                if (cur === floorStates.length - 1 || floorStates.slice(cur+1).lastIndexOf(true) === -1) {
                    dir = "down";
                }
            } else {
                if (cur === 0 || floorStates.slice(0, cur).indexOf(true) === -1) {
                    dir = "up";
                }
            }
            let next;
            if (dir === "up") {
                next = cur+1+floorStates.slice(cur+1).lastIndexOf(true);
            } else {
                next = floorStates.slice(0, cur).indexOf(true);
            }
            if (next !== -1) {
                elevator.goToFloor(next);
            } else {
                elevator.goToFloor(cur+ elevators == elevators[0] ? 1 : -1);
            }
        });
        elevator.on("floor_button_pressed", function(floorNum) {
            floorStates[floorNum] = true;
        });
        elevator.on("passing_floor", function(floorNum, direction) {
            if (floorStates[floorNum] || elevator.getPressedFloors().indexOf(floorNum) !== -1) {
                elevator.goToFloor(floorNum, true);
                floorStates[floorNum] = false;
            }
        });
        elevator.on("stopped_at_floor", function(floorNum) {
            floorStates[floorNum] = false;
        });
        for (let floor of floors) {
            floor.on("up_button_pressed", function() {
                floorStates[floor.floorNum()] = true;
            });
            floor.on("down_button_pressed", function() {
                floorStates[floor.floorNum()] = true;
            });
        }
        }
    },
    update: function(dt, elevators, floors) {
        // We normally don't need to do anything here
    }
}
