package glicko

import (
	"testing"
)

type player struct {
        *Rank
}

func newPlayer(r, rd float64) *player {
        p := new(player)
        p.Rank = &Rank{r, rd}
        return p
}

type players []*player

var testPlayers = []*player {
        newPlayer(1500, 200),
        newPlayer(1400, 30),
        newPlayer(1550, 100),
        newPlayer(1700, 300),
        newPlayer(1500, 350),
}

type varianceTest struct {
        p       *player
        in      Contests
        out     float64
}

func newContest(r, rd, outcome float64) *Contest {
        contest := new(Contest)
        contest.Rank = newPlayer(r, rd).Rank
        contest.Outcome = outcome
        return contest
}

type rdTest struct {
        in      *player
        out     float64
}

var rdTests = []*rdTest{
        &rdTest{testPlayers[0], 209.74803932337483},
        &rdTest{testPlayers[1], 69.95884504478329},
        &rdTest{testPlayers[2], 118.29725271535261},
}

var defaultOptions = &options{defaultMinRD, defaultMaxRD, defaultT, defaultC}

func TestRD(t *testing.T) {
	for _, ut := range rdTests {
                if newRD := rd(ut.in.Rank, defaultOptions); ut.out != newRD {
			t.Errorf("TestRD() = %+v, want %+v.", newRD, ut.out)
		}
	}
}

type gRDTest struct {
        in      float64
        out     float64
}

var gRDTests = []*gRDTest{
        &gRDTest{30, 0.9954980060779407},
        &gRDTest{100, 0.953148974234513},
        &gRDTest{300, 0.7242354637381511},
}

func TestGRD(t *testing.T) {
	for _, ut := range gRDTests {
                if gRD := g(ut.in); ut.out != gRD {
			t.Errorf("TestGRD() = %+v, want %+v.", gRD, ut.out)
		}
	}
}

type eTest struct {
        p       *player
        in      *Contest
        out     float64
}

var eTests = []*eTest{
        &eTest{testPlayers[0], newContest(testPlayers[1].R, testPlayers[1].RD, 1), 0.6369062639623734},
        &eTest{testPlayers[0], newContest(testPlayers[2].R, testPlayers[2].RD, 0), 0.43304012130623676},
        &eTest{testPlayers[0], newContest(testPlayers[3].R, testPlayers[3].RD, 0), 0.304672365112378},
}

func TestE(t *testing.T) {
	for _, ut := range eTests {
                if newE := e(ut.p.Rank, ut.in, defaultOptions); newE != ut.out {
			t.Errorf("TestE() = %+v, want %+v.", newE, ut.out)
		}
	}
}

type d2Test struct {
        p       *player
        in      Contests
        out     float64
}

var d2Tests = []*d2Test{
        &d2Test{
                testPlayers[0],
                Contests{
                        newContest(testPlayers[1].R, testPlayers[1].RD, 1),
                        newContest(testPlayers[2].R, testPlayers[2].RD, 0),
                        newContest(testPlayers[3].R, testPlayers[3].RD, 0),
                },
                55433.47321339554,
        },
}

func TestD2(t *testing.T) {
	for _, ut := range d2Tests {
                if newD2 := d2(ut.p.Rank, ut.in, defaultOptions); newD2 != ut.out {
			t.Errorf("TestD2() = %+v, want %+v.", newD2, ut.out)
		}
	}
}

type rPrimeTest struct {
        p       *player
        in      Contests
        out     float64
}

var rPrimeTests = []*rPrimeTest{
        &rPrimeTest{
                testPlayers[0],
                Contests{
                        newContest(testPlayers[1].R, testPlayers[1].RD, 1),
                        newContest(testPlayers[2].R, testPlayers[2].RD, 0),
                        newContest(testPlayers[3].R, testPlayers[3].RD, 0),
                },
                1461.9750469796434,
        },
}

func TestRPrime(t *testing.T) {
	for _, ut := range rPrimeTests {
                if newRPrime := rPrime(ut.p.Rank, ut.in, defaultOptions); newRPrime != ut.out {
			t.Errorf("TestRPrime() = %+v, want %+v.", newRPrime, ut.out)
		}
	}
}

type rdPrimeTest struct {
        p       *player
        in      Contests
        out     float64
}

var rdPrimeTests = []*rdPrimeTest{
        &rdPrimeTest{
                testPlayers[0],
                Contests{
                        newContest(testPlayers[1].R, testPlayers[1].RD, 1),
                        newContest(testPlayers[2].R, testPlayers[2].RD, 0),
                        newContest(testPlayers[3].R, testPlayers[3].RD, 0),
                },
                156.61387296904394,
        },
}

func TestRDPrime(t *testing.T) {
	for _, ut := range rdPrimeTests {
                if newRDPrime := rdPrime(ut.p.Rank, ut.in, defaultOptions); newRDPrime != ut.out {
			t.Errorf("TestRDPrime() = %+v, want %+v.", newRDPrime, ut.out)
		}
	}
}

type updateRatingTest struct {
        p       *player
        in      Contests
        out     *player
}

var updateRatingTests = []updateRatingTest{
        updateRatingTest{
                testPlayers[0],
                Contests{
                        newContest(testPlayers[1].R, testPlayers[1].RD, 1),
                        newContest(testPlayers[2].R, testPlayers[2].RD, 0),
                        newContest(testPlayers[3].R, testPlayers[3].RD, 0),
                },
                newPlayer(1461.9750469796434, 156.61387296904394),
        },
}

func TestUpdateRating(t *testing.T) {
	for _, ut := range updateRatingTests {
                newRating, err := UpdateRating(ut.p.Rank, ut.in)
                if err != nil {
			t.Errorf("UpdateRating() err = %v, expect nil.", err)
                }
                if !(newRating.RD == ut.out.Rank.RD && newRating.R == ut.out.Rank.R) {
			t.Errorf("UpdateRating() = %v %v, want %v %v.", newRating.R, newRating.RD, ut.out.R, ut.out.RD)
		}
	}
}

type updateDecayTest struct {
        p       *player
        in      Contests
        out     *player
}

var updateDecayTests = []updateDecayTest{
        updateDecayTest{
                testPlayers[1],
                Contests{},
                newPlayer(1400, 350),
        },
}

func TestDecayRating(t *testing.T) {
	for _, ut := range updateDecayTests {
                for i := 0; i < 31; i++ {
                        rating, err := UpdateRating(ut.p.Rank, ut.in)
                        if err != nil {
                                t.Errorf("UpdateRating() err = %v, expect nil.", err)
                        }
                        ut.p.R, ut.p.RD = rating.R, rating.RD
                }
                if !(ut.p.RD == ut.out.RD && ut.p.R == ut.out.R) {
			t.Errorf("DecayRating() = %v %v, want %v %v.", ut.p.R, ut.p.RD, ut.out.R, ut.out.RD)
		}
	}
}
