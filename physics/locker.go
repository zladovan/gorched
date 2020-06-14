package physics

// TimeLocker is tool for switching body lock after specified amount of seconds is passed.
// Time is running only when Update is called.
type TimeLocker struct {
	// BodyToRelock holds body which Lock attribute will be negated after time limit
	BodyToRelock *Body
	// RemainingSeconds holds number of seconds until Lock will be changed
	RemainingSeconds float64
	// relocked is flag for remembering that Lock was already changed
	relocked bool
}

// Update updates locker time and change body lock maximally once
func (t *TimeLocker) Update(dt float64) {
	if t.relocked {
		// already done
		return
	}

	// update remaining seconds
	if t.RemainingSeconds > 0 {
		t.RemainingSeconds -= dt
	}

	// change lock if it's time
	if t.RemainingSeconds <= 0 {
		t.BodyToRelock.Locked = !t.BodyToRelock.Locked
		t.relocked = true
	}
}
