// Implements the Glicko Rating System described at http://www.glicko.net/glicko.html.
package glicko

import (
	"errors"
	"math"
)

const (
	// A frequent player may end up with an RD approaching 0 which results in a rating that doesn't change despite improvement
	// A minimum RD helps prevent such stagnation.  The value defaults to 30 which is believed to be sufficient for most uses.
	defaultMinRD float64 = 30

	// RD shows uncertainty in the rating.  A rating for a person that has played games should be more certain than a new player.
	// MaxRD provides an upper limit and likely should be set to the RD for a new unranked player.
	// Default value is 350.
	defaultMaxRD float64 = 350

	// T represents the number of rating periods that have past since the present rating was last updated.
	// Defaults to 1.
	defaultT float64 = 1.0

	// C is a constant used to decay RD over time.
	// See, www.glicko.net for a discussion of its value.
	// Default is 63.2, which corresponds to a decay length of 30 periods when 350 is used for RD of unranked player.
	defaultC float64 = 63.2
)

//type Ranked interface {
//        // Rating Accessors
//        R() float64
//        SetR(float64)
//        // Rating Deviation Accessors
//        RD() float64
//        SetRD(float64)
//}

type Rank struct {
	R  float64
	RD float64
}

//func (this *Rank) R() float64 {
//        return this.r
//}
//
//func (this *Rank) SetR(r float64) {
//        this.r = r
//}
//
//func (this *Rank) RD() float64 {
//        return this.rd
//}
//
//func (this *Rank) SetRD(rd float64) {
//        this.rd = rd
//}

type options struct {
	minRD float64
	maxRD float64
	t     float64
	c     float64
}

// Contest records a single contest outcome against an opponent and that oppenents rating and rating deviation at the time.
type Contest struct {
	// 1 for win, 0 for loss, 0.5 for tie
	Outcome float64
	// Rating and Deviation of Opponent
	*Rank
}

func rd(player *Rank, o *options) float64 {
	return math.Min(math.Sqrt(player.RD*player.RD+o.c*o.c*o.t), o.maxRD)
}

func processOptions(os ...float64) (*options, error) {
	switch len(os) {
	case 0:
		return &options{defaultMinRD, defaultMaxRD, defaultT, defaultC}, nil
	case 4:
		return &options{os[0], os[1], os[2], os[3]}, nil
	}
	return nil, errors.New("Wrong number of options.  You must provide either 0 or 4 options.")
}

func rd2(player *Rank, o *options) float64 {
	newRD := rd(player, o)
	return newRD * newRD
}

// A collection of results for a single player during rating period
type Contests []*Contest

const q = 0.00575646273249 // Log(10) / 400
const q2 = q * q

func rPrime(rank *Rank, cs Contests, o *options) float64 {
	return math.Floor(rank.R + (q/((1/rd2(rank, o))+(1/d2(rank, cs, o))))*gsSum(rank, cs, o))
}

func d2(rank *Rank, cs Contests, o *options) float64 {
	sum := 0.0
	for _, c := range cs {
		newRD := rd(c.Rank, o)
		gRD := g(newRD)
		sum += gRD * gRD * e1e(rank, c, o)
	}
	return 1 / (q2 * sum)
}

func g(rd float64) float64 {
	return 1 / math.Sqrt(1+(3*q2*math.Pow(rd, 2))/math.Pow(math.Pi, 2))
}

func e(rank *Rank, contest *Contest, o *options) float64 {
	return 1 / (1 + math.Pow(10, g(rd(contest.Rank, o))*(rank.R-contest.R)/-400))
}

func e1e(rank *Rank, contest *Contest, o *options) float64 {
	ec := e(rank, contest, o)
	return ec * (1 - ec)
}

func gsSum(rank *Rank, cs Contests, o *options) (sum float64) {
	for _, c := range cs {
		sum += g(rd(c.Rank, o)) * (c.Outcome - e(rank, c, o))
	}
	return sum
}

func rdPrime(rank *Rank, cs Contests, o *options) float64 {
	return math.Max(math.Floor(math.Sqrt(1/((1/rd2(rank, o))+(1/d2(rank, cs, o))))), o.minRD)
}

// Returns updated rating and updated rating deviation based on the provided contest results.
// If len(options) == 0, default values for minRD, maxRD, t, and c are used.
// Otherwise, must provide minRD, maxRD, t, and c.
func UpdateRating(rank *Rank, cs Contests, options ...float64) (*Rank, error) {
	o, err := processOptions(options...)
	switch {
	case err != nil:
		return nil, err
	case len(cs) == 0:
		return &Rank{rank.R, rd(rank, o)}, nil
	}
	return &Rank{rPrime(rank, cs, o), rdPrime(rank, cs, o)}, nil
}

func ConfidenceInterval(rank *Rank) (float64, float64) {
	twiceRd := 2.0 * rank.RD
	return rank.R - twiceRd, rank.R + twiceRd
}
