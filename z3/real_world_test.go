// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package z3

import (
	"testing"
)

// TestRabbitsAndPheasants solves the classic riddle:
// 9 animals, rabbits and pheasants are playing on the grass.
// We can see 24 legs. How many rabbits and pheasants are there?
// Rabbits have 4 legs, pheasants have 2 legs.
// Based on: https://www.keiruaprod.fr/blog/2021/05/09/z3-samples.html
func TestRabbitsAndPheasants(t *testing.T) {
	ctx := NewContext(nil)
	solver := NewSolver(ctx)

	rabbits := ctx.IntConst("rabbits")
	pheasants := ctx.IntConst("pheasants")

	zero := ctx.Int(0)
	nine := ctx.Int(9)
	twentyFour := ctx.Int(24)
	two := ctx.Int(2)
	four := ctx.Int(4)

	// Total animals: rabbits + pheasants == 9
	solver.Assert(rabbits.Add(pheasants).Eq(nine))
	// Total legs: 4*rabbits + 2*pheasants == 24
	solver.Assert(rabbits.Mul(four).Add(pheasants.Mul(two)).Eq(twentyFour))
	// Both counts must be non-negative
	solver.Assert(rabbits.GE(zero))
	solver.Assert(pheasants.GE(zero))

	sat, err := solver.Check()
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if !sat {
		t.Fatal("expected satisfiable")
	}

	model := solver.Model()
	rabbitsVal, _, ok := model.EvalAsInt64(rabbits, true)
	if !ok {
		t.Fatal("could not evaluate rabbits")
	}
	pheasantsVal, _, ok := model.EvalAsInt64(pheasants, true)
	if !ok {
		t.Fatal("could not evaluate pheasants")
	}

	t.Logf("Rabbits: %d, Pheasants: %d", rabbitsVal, pheasantsVal)

	// Verify the solution
	if rabbitsVal != 3 {
		t.Fatalf("expected 3 rabbits, got %d", rabbitsVal)
	}
	if pheasantsVal != 6 {
		t.Fatalf("expected 6 pheasants, got %d", pheasantsVal)
	}
}

// TestRabbitsAndPheasantsWithOr solves a variant:
// 9 animals, rabbits and pheasants. We can see 24 or 27 legs.
// How many rabbits and pheasants are there?
// Only 24 legs gives an integer solution.
func TestRabbitsAndPheasantsWithOr(t *testing.T) {
	ctx := NewContext(nil)
	solver := NewSolver(ctx)

	rabbits := ctx.IntConst("rabbits")
	pheasants := ctx.IntConst("pheasants")

	zero := ctx.Int(0)
	nine := ctx.Int(9)
	twentyFour := ctx.Int(24)
	twentySeven := ctx.Int(27)
	two := ctx.Int(2)
	four := ctx.Int(4)

	// Total animals: rabbits + pheasants == 9
	solver.Assert(rabbits.Add(pheasants).Eq(nine))
	// Total legs: 4*rabbits + 2*pheasants == 24 OR == 27
	totalLegs := rabbits.Mul(four).Add(pheasants.Mul(two))
	solver.Assert(totalLegs.Eq(twentyFour).Or(totalLegs.Eq(twentySeven)))
	// Both counts must be non-negative
	solver.Assert(rabbits.GE(zero))
	solver.Assert(pheasants.GE(zero))

	sat, err := solver.Check()
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if !sat {
		t.Fatal("expected satisfiable")
	}

	model := solver.Model()
	rabbitsVal, _, ok := model.EvalAsInt64(rabbits, true)
	if !ok {
		t.Fatal("could not evaluate rabbits")
	}
	pheasantsVal, _, ok := model.EvalAsInt64(pheasants, true)
	if !ok {
		t.Fatal("could not evaluate pheasants")
	}

	t.Logf("Rabbits: %d, Pheasants: %d", rabbitsVal, pheasantsVal)

	// Verify it's a valid solution (24 legs case)
	totalLegsVal := 4*rabbitsVal + 2*pheasantsVal
	if totalLegsVal != 24 && totalLegsVal != 27 {
		t.Fatalf("expected 24 or 27 legs, got %d", totalLegsVal)
	}
	if rabbitsVal+pheasantsVal != 9 {
		t.Fatalf("expected 9 animals, got %d", rabbitsVal+pheasantsVal)
	}
}

// TestXKCD287 solves the XKCD 287 "NP-Complete" problem:
// Find combinations of appetizers that sum to exactly $15.05
// Prices: Mixed fruits $2.15, French Fries $2.75, Side Salad $3.35,
//
//	Hot Wings $3.55, Mozzarella Sticks $4.20, Sampler Plate $5.80
//
// Based on: https://www.keiruaprod.fr/blog/2021/05/09/z3-samples.html
func TestXKCD287(t *testing.T) {
	ctx := NewContext(nil)
	solver := NewSolver(ctx)

	// Prices in cents to avoid floating point issues
	prices := []int64{215, 275, 335, 355, 420, 580}
	appetizers := []string{
		"Mixed fruits",
		"French Fries",
		"Side Salad",
		"Hot Wings",
		"Mozzarella Sticks",
		"Sampler Plate",
	}
	total := int64(1505)

	// Create quantity variables for each appetizer
	quantities := make([]Int, len(appetizers))
	for i := range appetizers {
		quantities[i] = ctx.IntConst(appetizers[i])
		// Quantities must be between 0 and 10
		solver.Assert(quantities[i].GE(ctx.Int(0)))
		solver.Assert(quantities[i].LE(ctx.Int(10)))
	}

	// Sum of (quantity * price) must equal total
	var sumTerms []Int
	for i, price := range prices {
		priceVal := ctx.Int64(price)
		sumTerms = append(sumTerms, quantities[i].Mul(priceVal))
	}

	// Build the sum constraint
	totalSum := sumTerms[0]
	for i := 1; i < len(sumTerms); i++ {
		totalSum = totalSum.Add(sumTerms[i])
	}
	solver.Assert(totalSum.Eq(ctx.Int64(total)))

	// Find all solutions
	solutions := 0
	for {
		sat, err := solver.Check()
		if err != nil {
			t.Fatalf("error: %s", err)
		}
		if !sat {
			break
		}

		solutions++
		model := solver.Model()

		t.Logf("Solution %d:", solutions)
		var blocking []Bool
		for i, name := range appetizers {
			qtyVal, _, ok := model.EvalAsInt64(quantities[i], true)
			if !ok {
				t.Fatalf("could not evaluate %s", name)
			}
			if qtyVal > 0 {
				t.Logf("  %d x %s = $%.2f", qtyVal, name, float64(qtyVal)*float64(prices[i])/100.0)
			}
			// Add constraint to exclude this solution
			blocking = append(blocking, quantities[i].NE(ctx.Int64(qtyVal)))
		}

		// Add constraint to find different solutions
		if len(blocking) > 0 {
			or := blocking[0]
			for i := 1; i < len(blocking); i++ {
				or = or.Or(blocking[i])
			}
			solver.Assert(or)
		}

		if solutions >= 10 {
			t.Log("Stopping after 10 solutions")
			break
		}
	}

	if solutions < 1 {
		t.Fatal("expected at least one solution")
	}
	t.Logf("Found %d solution(s)", solutions)
}

// TestEinsteinRiddle solves the famous Zebra Puzzle / Einstein's Riddle:
// - There are five houses.
// - The Englishman lives in the red house.
// - The Spaniard owns the dog.
// - Coffee is drunk in the green house.
// - The Ukrainian drinks tea.
// - The green house is immediately to the right of the ivory house.
// - The Old Gold smoker owns snails.
// - Kools are smoked in the yellow house.
// - Milk is drunk in the middle house.
// - The Norwegian lives in the first house.
// - The man who smokes Chesterfields lives next to the man with the fox.
// - Kools are smoked in the house next to the house where the horse is kept.
// - The Lucky Strike smoker drinks orange juice.
// - The Japanese smokes Parliaments.
// - The Norwegian lives next to the blue house.
// Question: Who drinks water? Who owns the zebra?
// Based on: https://www.keiruaprod.fr/blog/2021/05/09/z3-samples.html
func TestEinsteinRiddle(t *testing.T) {
	ctx := NewContext(nil)
	solver := NewSolver(ctx)

	// Each variable represents the house number (1-5) where that person/item is
	// Nationalities
	englishman := ctx.IntConst("Englishman")
	spaniard := ctx.IntConst("Spaniard")
	ukrainian := ctx.IntConst("Ukrainian")
	norwegian := ctx.IntConst("Norwegian")
	japanese := ctx.IntConst("Japanese")
	nationalities := []Int{englishman, spaniard, ukrainian, norwegian, japanese}

	// Cigarettes
	parliaments := ctx.IntConst("Parliaments")
	kools := ctx.IntConst("Kools")
	luckyStrike := ctx.IntConst("LuckyStrike")
	oldGold := ctx.IntConst("OldGold")
	chesterfields := ctx.IntConst("Chesterfields")
	cigarettes := []Int{parliaments, kools, luckyStrike, oldGold, chesterfields}

	// Animals
	fox := ctx.IntConst("Fox")
	horse := ctx.IntConst("Horse")
	zebra := ctx.IntConst("Zebra")
	dog := ctx.IntConst("Dog")
	snails := ctx.IntConst("Snails")
	animals := []Int{fox, horse, zebra, dog, snails}

	// Drinks
	coffee := ctx.IntConst("Coffee")
	milk := ctx.IntConst("Milk")
	orangeJuice := ctx.IntConst("OrangeJuice")
	tea := ctx.IntConst("Tea")
	water := ctx.IntConst("Water")
	drinks := []Int{coffee, milk, orangeJuice, tea, water}

	// Colors
	red := ctx.IntConst("Red")
	green := ctx.IntConst("Green")
	ivory := ctx.IntConst("Ivory")
	blue := ctx.IntConst("Blue")
	yellow := ctx.IntConst("Yellow")
	colors := []Int{red, green, ivory, blue, yellow}

	allGroups := [][]Int{nationalities, cigarettes, animals, drinks, colors}

	one := ctx.Int(1)
	three := ctx.Int(3)
	five := ctx.Int(5)

	// Helper function: neighbor constraint (|a - b| == 1)
	neighbor := func(a, b Int) Bool {
		diff := a.Sub(b)
		return diff.Eq(one).Or(diff.Eq(ctx.Int(-1)))
	}

	// Constraints: Each category has distinct values 1-5
	for _, group := range allGroups {
		for _, v := range group {
			solver.Assert(v.GE(one))
			solver.Assert(v.LE(five))
		}
		// All different within group
		for i := 0; i < len(group); i++ {
			for j := i + 1; j < len(group); j++ {
				solver.Assert(group[i].NE(group[j]))
			}
		}
	}

	// Clues:
	// 1. The Englishman lives in the red house.
	solver.Assert(englishman.Eq(red))
	// 2. The Spaniard owns the dog.
	solver.Assert(spaniard.Eq(dog))
	// 3. Coffee is drunk in the green house.
	solver.Assert(coffee.Eq(green))
	// 4. The Ukrainian drinks tea.
	solver.Assert(ukrainian.Eq(tea))
	// 5. The green house is immediately to the right of the ivory house.
	solver.Assert(green.Eq(ivory.Add(one)))
	// 6. The Old Gold smoker owns snails.
	solver.Assert(oldGold.Eq(snails))
	// 7. Kools are smoked in the yellow house.
	solver.Assert(kools.Eq(yellow))
	// 8. Milk is drunk in the middle house.
	solver.Assert(milk.Eq(three))
	// 9. The Norwegian lives in the first house.
	solver.Assert(norwegian.Eq(one))
	// 10. The man who smokes Chesterfields lives next to the man with the fox.
	solver.Assert(neighbor(chesterfields, fox))
	// 11. Kools are smoked in the house next to the house where the horse is kept.
	solver.Assert(neighbor(kools, horse))
	// 12. The Lucky Strike smoker drinks orange juice.
	solver.Assert(luckyStrike.Eq(orangeJuice))
	// 13. The Japanese smokes Parliaments.
	solver.Assert(japanese.Eq(parliaments))
	// 14. The Norwegian lives next to the blue house.
	solver.Assert(neighbor(norwegian, blue))

	sat, err := solver.Check()
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if !sat {
		t.Fatal("expected satisfiable")
	}

	model := solver.Model()

	// Find who drinks water
	waterHouse, _, _ := model.EvalAsInt64(water, true)

	// Find who owns the zebra
	zebraHouse, _, _ := model.EvalAsInt64(zebra, true)

	// Find the nationalities in those houses
	var waterDrinker, zebraOwner string
	nationalityNames := []string{"Englishman", "Spaniard", "Ukrainian", "Norwegian", "Japanese"}
	for i, nat := range nationalities {
		house, _, _ := model.EvalAsInt64(nat, true)
		if house == waterHouse {
			waterDrinker = nationalityNames[i]
		}
		if house == zebraHouse {
			zebraOwner = nationalityNames[i]
		}
	}

	t.Logf("Who drinks water? %s", waterDrinker)
	t.Logf("Who owns the zebra? %s", zebraOwner)

	// The known answer: Norwegian drinks water, Japanese owns the zebra
	if waterDrinker != "Norwegian" {
		t.Fatalf("expected Norwegian drinks water, got %s", waterDrinker)
	}
	if zebraOwner != "Japanese" {
		t.Fatalf("expected Japanese owns zebra, got %s", zebraOwner)
	}
}

// TestSkisAssignment solves an optimization problem:
// Assign skis to skiers minimizing the total disparity between ski sizes and skier heights.
// Ski sizes: 1, 2, 5, 7, 13, 21
// Skier heights: 3, 4, 7, 11, 18
// Based on: https://www.keiruaprod.fr/blog/2021/05/09/z3-samples.html
func TestSkisAssignment(t *testing.T) {
	ctx := NewContext(nil)
	opt := NewOptimize(ctx)

	skiSizes := []int64{1, 2, 5, 7, 13, 21}
	skierHeights := []int64{3, 4, 7, 11, 18}

	zero := ctx.Int(0)
	numSkis := ctx.Int(len(skiSizes))

	// Create assignment variables: assignments[i] = ski index for skier i
	assignments := make([]Int, len(skierHeights))
	for i := range skierHeights {
		assignments[i] = ctx.IntConst("ski_for_skier_" + string(rune('0'+i)))
		// Each assignment must be a valid ski index
		opt.Assert(assignments[i].GE(zero))
		opt.Assert(assignments[i].LT(numSkis))
	}

	// All assignments must be different (each skier gets a different ski)
	for i := 0; i < len(assignments); i++ {
		for j := i + 1; j < len(assignments); j++ {
			opt.Assert(assignments[i].NE(assignments[j]))
		}
	}

	// Calculate total disparity to minimize
	// We need to compute sum of |skiSize[assignment[i]] - skierHeight[i]|
	// Using a helper approach: for each skier, we use a disparity variable

	var disparities []Int
	for i, height := range skierHeights {
		disparity := ctx.IntConst("disparity_" + string(rune('0'+i)))
		disparities = append(disparities, disparity)

		// disparity[i] = |skiSize[assignment[i]] - height|
		// We need to encode: for each possible ski j, if assignment[i] == j then disparity == |skiSize[j] - height|
		for j, skiSize := range skiSizes {
			diff := skiSize - height
			if diff < 0 {
				diff = -diff
			}
			diffVal := ctx.Int64(diff)
			jVal := ctx.Int(j)
			// If assignment == j, then disparity == |diff|
			opt.Assert(assignments[i].Eq(jVal).Implies(disparity.Eq(diffVal)))
		}
		// Disparity must be non-negative
		opt.Assert(disparity.GE(zero))
	}

	// Total disparity to minimize
	totalDisparity := disparities[0]
	for i := 1; i < len(disparities); i++ {
		totalDisparity = totalDisparity.Add(disparities[i])
	}

	obj := opt.Minimize(totalDisparity)

	sat, err := opt.Check()
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if !sat {
		t.Fatal("expected satisfiable")
	}

	model := opt.Model()
	t.Logf("Minimum total disparity: %s", obj.Lower())

	for i, height := range skierHeights {
		assignVal, _, _ := model.EvalAsInt64(assignments[i], true)
		skiSize := skiSizes[assignVal]
		dispVal, _, _ := model.EvalAsInt64(disparities[i], true)
		t.Logf("Skier %d (height %d) gets ski of size %d (disparity: %d)", i, height, skiSize, dispVal)
	}
}

// TestOrganizeYourDay solves a scheduling problem:
// Your day starts at 9 and finishes at 17.
// Tasks: work (4 hours), mail (1 hour), bank (2 hours), shopping (1 hour)
// Constraints:
// - One task has to be finished before another starts (no overlap)
// - Send the mail before going to work
// - Go to the bank before shopping
// - Start work after 11
// Based on: https://www.keiruaprod.fr/blog/2021/05/09/z3-samples.html
func TestOrganizeYourDay(t *testing.T) {
	ctx := NewContext(nil)
	solver := NewSolver(ctx)

	// Start times for each task
	workStart := ctx.IntConst("work_start")
	mailStart := ctx.IntConst("mail_start")
	bankStart := ctx.IntConst("bank_start")
	shoppingStart := ctx.IntConst("shopping_start")

	tasks := []Int{workStart, mailStart, bankStart, shoppingStart}
	durations := []int64{4, 1, 2, 1}
	taskNames := []string{"work", "mail", "bank", "shopping"}

	nine := ctx.Int(9)
	eleven := ctx.Int(11)
	seventeen := ctx.Int(17)

	// Each task must start after 9 and finish by 17
	for i := range tasks {
		durationVal := ctx.Int64(durations[i])
		solver.Assert(tasks[i].GE(nine))
		solver.Assert(tasks[i].Add(durationVal).LE(seventeen))
	}

	// No overlap: for any two tasks, one must finish before the other starts
	for i := 0; i < len(tasks); i++ {
		for j := i + 1; j < len(tasks); j++ {
			duration1 := ctx.Int64(durations[i])
			duration2 := ctx.Int64(durations[j])
			// task1 finishes before task2 starts OR task2 finishes before task1 starts
			solver.Assert(
				tasks[i].Add(duration1).LE(tasks[j]).Or(
					tasks[j].Add(duration2).LE(tasks[i]),
				),
			)
		}
	}

	// Additional constraints:
	// - Start work after 11
	solver.Assert(workStart.GE(eleven))
	// - Send the mail before going to work
	mailDuration := ctx.Int(1)
	solver.Assert(mailStart.Add(mailDuration).LE(workStart))
	// - Go to the bank before shopping
	bankDuration := ctx.Int(2)
	solver.Assert(bankStart.Add(bankDuration).LE(shoppingStart))

	sat, err := solver.Check()
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if !sat {
		t.Fatal("expected satisfiable")
	}

	model := solver.Model()

	t.Log("Schedule:")
	for i := range tasks {
		startVal, _, _ := model.EvalAsInt64(tasks[i], true)
		duration := durations[i]
		name := taskNames[i]
		t.Logf("  %s: %d:00 - %d:00 (%d hour(s))", name, startVal, startVal+duration, duration)
	}

	// Verify constraints
	workVal, _, _ := model.EvalAsInt64(workStart, true)
	mailVal, _, _ := model.EvalAsInt64(mailStart, true)
	bankVal, _, _ := model.EvalAsInt64(bankStart, true)
	shoppingVal, _, _ := model.EvalAsInt64(shoppingStart, true)

	if workVal < 11 {
		t.Fatalf("work should start after 11, got %d", workVal)
	}
	if mailVal+1 > workVal {
		t.Fatalf("mail should finish before work starts")
	}
	if bankVal+2 > shoppingVal {
		t.Fatalf("bank should finish before shopping starts")
	}
}

// TestSudoku solves a classic 9x9 Sudoku puzzle using Z3.
// This demonstrates using integer variables with domain constraints.
func TestSudoku(t *testing.T) {
	ctx := NewContext(nil)
	solver := NewSolver(ctx)

	one := ctx.Int(1)
	nine := ctx.Int(9)

	// Create 9x9 grid of integer variables
	cells := make([][]Int, 9)
	for i := 0; i < 9; i++ {
		cells[i] = make([]Int, 9)
		for j := 0; j < 9; j++ {
			cells[i][j] = ctx.IntConst("cell_" + string(rune('0'+i)) + "_" + string(rune('0'+j)))
			// Each cell is between 1 and 9
			solver.Assert(cells[i][j].GE(one))
			solver.Assert(cells[i][j].LE(nine))
		}
	}

	// Each row has distinct values
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			for k := j + 1; k < 9; k++ {
				solver.Assert(cells[i][j].NE(cells[i][k]))
			}
		}
	}

	// Each column has distinct values
	for j := 0; j < 9; j++ {
		for i := 0; i < 9; i++ {
			for k := i + 1; k < 9; k++ {
				solver.Assert(cells[i][j].NE(cells[k][j]))
			}
		}
	}

	// Each 3x3 box has distinct values
	for boxRow := 0; boxRow < 3; boxRow++ {
		for boxCol := 0; boxCol < 3; boxCol++ {
			var boxCells []Int
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					boxCells = append(boxCells, cells[boxRow*3+i][boxCol*3+j])
				}
			}
			for i := 0; i < len(boxCells); i++ {
				for j := i + 1; j < len(boxCells); j++ {
					solver.Assert(boxCells[i].NE(boxCells[j]))
				}
			}
		}
	}

	// Sample puzzle (0 means empty)
	// This is an easy puzzle for testing
	puzzle := [][]int{
		{5, 3, 0, 0, 7, 0, 0, 0, 0},
		{6, 0, 0, 1, 9, 5, 0, 0, 0},
		{0, 9, 8, 0, 0, 0, 0, 6, 0},
		{8, 0, 0, 0, 6, 0, 0, 0, 3},
		{4, 0, 0, 8, 0, 3, 0, 0, 1},
		{7, 0, 0, 0, 2, 0, 0, 0, 6},
		{0, 6, 0, 0, 0, 0, 2, 8, 0},
		{0, 0, 0, 4, 1, 9, 0, 0, 5},
		{0, 0, 0, 0, 8, 0, 0, 7, 9},
	}

	// Add initial constraints for given values
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if puzzle[i][j] != 0 {
				solver.Assert(cells[i][j].Eq(ctx.Int(puzzle[i][j])))
			}
		}
	}

	sat, err := solver.Check()
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if !sat {
		t.Fatal("expected satisfiable")
	}

	model := solver.Model()

	// Print solution
	t.Log("Sudoku solution:")
	for i := 0; i < 9; i++ {
		row := ""
		for j := 0; j < 9; j++ {
			val, _, _ := model.EvalAsInt64(cells[i][j], true)
			row += string(rune('0' + val))
			if j == 2 || j == 5 {
				row += "|"
			}
		}
		t.Log(row)
		if i == 2 || i == 5 {
			t.Log("---+---+---")
		}
	}
}

// TestNQueens solves the N-Queens problem:
// Place N queens on an NxN chessboard so that no two queens attack each other.
func TestNQueens(t *testing.T) {
	ctx := NewContext(nil)
	solver := NewSolver(ctx)

	const n = 8

	// queens[i] represents the column position of the queen in row i
	queens := make([]Int, n)
	for i := 0; i < n; i++ {
		queens[i] = ctx.IntConst("queen_" + string(rune('0'+i)))
		// Each queen is in a valid column (0 to n-1)
		solver.Assert(queens[i].GE(ctx.Int(0)))
		solver.Assert(queens[i].LT(ctx.Int(n)))
	}

	// No two queens in the same column
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			solver.Assert(queens[i].NE(queens[j]))
		}
	}

	// No two queens on the same diagonal
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			// |queens[i] - queens[j]| != |i - j|
			diff := j - i
			diffVal := ctx.Int(diff)
			negDiffVal := ctx.Int(-diff)
			// queens[j] - queens[i] != diff AND queens[j] - queens[i] != -diff
			colDiff := queens[j].Sub(queens[i])
			solver.Assert(colDiff.NE(diffVal))
			solver.Assert(colDiff.NE(negDiffVal))
		}
	}

	sat, err := solver.Check()
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if !sat {
		t.Fatal("expected satisfiable")
	}

	model := solver.Model()

	t.Logf("%d-Queens solution:", n)
	for i := 0; i < n; i++ {
		col, _, _ := model.EvalAsInt64(queens[i], true)
		row := ""
		for j := 0; j < n; j++ {
			if int64(j) == col {
				row += "Q "
			} else {
				row += ". "
			}
		}
		t.Log(row)
	}
}

// TestMagicSquare solves a 3x3 magic square:
// Each row, column, and diagonal sums to the same value (magic constant = 15).
func TestMagicSquare(t *testing.T) {
	ctx := NewContext(nil)
	solver := NewSolver(ctx)

	n := 3
	magicSum := 15 // For 3x3, the magic constant is 15

	one := ctx.Int(1)
	nine := ctx.Int(n * n)
	target := ctx.Int(magicSum)

	// Create n x n grid
	cells := make([][]Int, n)
	var allCells []Int
	for i := 0; i < n; i++ {
		cells[i] = make([]Int, n)
		for j := 0; j < n; j++ {
			cells[i][j] = ctx.IntConst("m_" + string(rune('0'+i)) + "_" + string(rune('0'+j)))
			// Each cell contains 1 to n*n
			solver.Assert(cells[i][j].GE(one))
			solver.Assert(cells[i][j].LE(nine))
			allCells = append(allCells, cells[i][j])
		}
	}

	// All cells have distinct values
	for i := 0; i < len(allCells); i++ {
		for j := i + 1; j < len(allCells); j++ {
			solver.Assert(allCells[i].NE(allCells[j]))
		}
	}

	// Row sums
	for i := 0; i < n; i++ {
		rowSum := cells[i][0]
		for j := 1; j < n; j++ {
			rowSum = rowSum.Add(cells[i][j])
		}
		solver.Assert(rowSum.Eq(target))
	}

	// Column sums
	for j := 0; j < n; j++ {
		colSum := cells[0][j]
		for i := 1; i < n; i++ {
			colSum = colSum.Add(cells[i][j])
		}
		solver.Assert(colSum.Eq(target))
	}

	// Main diagonal
	diagSum := cells[0][0]
	for i := 1; i < n; i++ {
		diagSum = diagSum.Add(cells[i][i])
	}
	solver.Assert(diagSum.Eq(target))

	// Anti-diagonal
	antiDiagSum := cells[0][n-1]
	for i := 1; i < n; i++ {
		antiDiagSum = antiDiagSum.Add(cells[i][n-1-i])
	}
	solver.Assert(antiDiagSum.Eq(target))

	sat, err := solver.Check()
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if !sat {
		t.Fatal("expected satisfiable")
	}

	model := solver.Model()

	t.Log("Magic Square solution:")
	for i := 0; i < n; i++ {
		row := ""
		for j := 0; j < n; j++ {
			val, _, _ := model.EvalAsInt64(cells[i][j], true)
			row += string(rune('0'+val)) + " "
		}
		t.Log(row)
	}
}

// TestGraphColoring solves a graph coloring problem:
// Color a graph with at most k colors such that no two adjacent vertices share a color.
func TestGraphColoring(t *testing.T) {
	ctx := NewContext(nil)
	solver := NewSolver(ctx)

	// Simple graph: 4 vertices forming a square with one diagonal
	// Edges: 0-1, 1-2, 2-3, 3-0, 0-2 (diagonal)
	numVertices := 4
	edges := [][2]int{{0, 1}, {1, 2}, {2, 3}, {3, 0}, {0, 2}}
	numColors := 3 // Try with 3 colors

	zero := ctx.Int(0)
	maxColor := ctx.Int(numColors - 1)

	// Create color variable for each vertex
	colors := make([]Int, numVertices)
	for i := 0; i < numVertices; i++ {
		colors[i] = ctx.IntConst("color_" + string(rune('0'+i)))
		solver.Assert(colors[i].GE(zero))
		solver.Assert(colors[i].LE(maxColor))
	}

	// Adjacent vertices must have different colors
	for _, edge := range edges {
		solver.Assert(colors[edge[0]].NE(colors[edge[1]]))
	}

	sat, err := solver.Check()
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if !sat {
		t.Fatal("expected satisfiable")
	}

	model := solver.Model()

	colorNames := []string{"Red", "Green", "Blue"}
	t.Log("Graph coloring solution:")
	for i := 0; i < numVertices; i++ {
		colorVal, _, _ := model.EvalAsInt64(colors[i], true)
		t.Logf("  Vertex %d: %s", i, colorNames[colorVal])
	}
}

// TestKnapsack solves a simple 0/1 knapsack problem:
// Maximize the value of items that fit in a knapsack with limited capacity.
func TestKnapsack(t *testing.T) {
	ctx := NewContext(nil)
	opt := NewOptimize(ctx)

	// Items: (weight, value)
	items := []struct {
		name   string
		weight int64
		value  int64
	}{
		{"laptop", 3, 10},
		{"camera", 2, 8},
		{"phone", 1, 5},
		{"book", 2, 3},
		{"snacks", 1, 2},
		{"headphones", 1, 4},
	}
	capacity := 6

	zero := ctx.Int(0)
	one := ctx.Int(1)
	capVal := ctx.Int(capacity)

	// Binary decision variables: take[i] = 0 or 1
	take := make([]Int, len(items))
	for i := range items {
		take[i] = ctx.IntConst("take_" + items[i].name)
		opt.Assert(take[i].GE(zero))
		opt.Assert(take[i].LE(one))
	}

	// Total weight constraint
	weightSum := take[0].Mul(ctx.Int64(items[0].weight))
	for i := 1; i < len(items); i++ {
		weightSum = weightSum.Add(take[i].Mul(ctx.Int64(items[i].weight)))
	}
	opt.Assert(weightSum.LE(capVal))

	// Maximize total value
	valueSum := take[0].Mul(ctx.Int64(items[0].value))
	for i := 1; i < len(items); i++ {
		valueSum = valueSum.Add(take[i].Mul(ctx.Int64(items[i].value)))
	}
	obj := opt.Maximize(valueSum)

	sat, err := opt.Check()
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if !sat {
		t.Fatal("expected satisfiable")
	}

	model := opt.Model()
	t.Logf("Maximum value: %s", obj.Upper())

	totalWeight := int64(0)
	t.Log("Selected items:")
	for i := range items {
		takeVal, _, _ := model.EvalAsInt64(take[i], true)
		if takeVal == 1 {
			t.Logf("  %s (weight: %d, value: %d)", items[i].name, items[i].weight, items[i].value)
			totalWeight += items[i].weight
		}
	}
	t.Logf("Total weight: %d / %d", totalWeight, capacity)
}

// TestSendMoreMoney solves the classic cryptarithmetic puzzle:
// SEND + MORE = MONEY where each letter represents a unique digit.
func TestSendMoreMoney(t *testing.T) {
	ctx := NewContext(nil)
	solver := NewSolver(ctx)

	zero := ctx.Int(0)
	nine := ctx.Int(9)
	one := ctx.Int(1)

	// Create a variable for each letter
	s := ctx.IntConst("S")
	e := ctx.IntConst("E")
	n := ctx.IntConst("N")
	d := ctx.IntConst("D")
	m := ctx.IntConst("M")
	o := ctx.IntConst("O")
	r := ctx.IntConst("R")
	y := ctx.IntConst("Y")

	letters := []Int{s, e, n, d, m, o, r, y}

	// Each letter is a digit 0-9
	for _, letter := range letters {
		solver.Assert(letter.GE(zero))
		solver.Assert(letter.LE(nine))
	}

	// All letters are different
	for i := 0; i < len(letters); i++ {
		for j := i + 1; j < len(letters); j++ {
			solver.Assert(letters[i].NE(letters[j]))
		}
	}

	// Leading digits cannot be zero
	solver.Assert(s.GE(one))
	solver.Assert(m.GE(one))

	// SEND + MORE = MONEY
	// SEND = 1000*S + 100*E + 10*N + D
	// MORE = 1000*M + 100*O + 10*R + E
	// MONEY = 10000*M + 1000*O + 100*N + 10*E + Y
	thousand := ctx.Int(1000)
	hundred := ctx.Int(100)
	ten := ctx.Int(10)
	tenThousand := ctx.Int(10000)

	send := s.Mul(thousand).Add(e.Mul(hundred)).Add(n.Mul(ten)).Add(d)
	more := m.Mul(thousand).Add(o.Mul(hundred)).Add(r.Mul(ten)).Add(e)
	money := m.Mul(tenThousand).Add(o.Mul(thousand)).Add(n.Mul(hundred)).Add(e.Mul(ten)).Add(y)

	solver.Assert(send.Add(more).Eq(money))

	sat, err := solver.Check()
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if !sat {
		t.Fatal("expected satisfiable")
	}

	model := solver.Model()

	sVal, _, _ := model.EvalAsInt64(s, true)
	eVal, _, _ := model.EvalAsInt64(e, true)
	nVal, _, _ := model.EvalAsInt64(n, true)
	dVal, _, _ := model.EvalAsInt64(d, true)
	mVal, _, _ := model.EvalAsInt64(m, true)
	oVal, _, _ := model.EvalAsInt64(o, true)
	rVal, _, _ := model.EvalAsInt64(r, true)
	yVal, _, _ := model.EvalAsInt64(y, true)

	sendNum := sVal*1000 + eVal*100 + nVal*10 + dVal
	moreNum := mVal*1000 + oVal*100 + rVal*10 + eVal
	moneyNum := mVal*10000 + oVal*1000 + nVal*100 + eVal*10 + yVal

	t.Logf("S=%d, E=%d, N=%d, D=%d, M=%d, O=%d, R=%d, Y=%d", sVal, eVal, nVal, dVal, mVal, oVal, rVal, yVal)
	t.Logf("SEND=%d + MORE=%d = MONEY=%d", sendNum, moreNum, moneyNum)

	if sendNum+moreNum != moneyNum {
		t.Fatalf("expected SEND + MORE = MONEY, got %d + %d = %d", sendNum, moreNum, moneyNum)
	}
}
