package question

type Question struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	// the true answer, 0-100 (%)
	Answer int `json:"answer"`
}

// question bank, each answer is a percentage (0-100)
// the game picks up to game.Rounds of them at random each game
var Bank = []Question{
	{ID: "q1", Text: "What percentage of people have ridden a motorcycle?", Answer: 25},
	{ID: "q2", Text: "What percentage of Earth's surface is covered by water?", Answer: 71},
	{ID: "q3", Text: "What percentage of Earth's atmosphere is nitrogen?", Answer: 78},
	{ID: "q4", Text: "What percentage of Earth's atmosphere is oxygen?", Answer: 21},
	{ID: "q5", Text: "What percentage of Earth's water is fresh water?", Answer: 3},
	{ID: "q6", Text: "Roughly what percentage of an adult human body is water?", Answer: 60},

	// space
	{ID: "q7", Text: "What percentage of the solar system's mass is in the Sun?", Answer: 99},
	{ID: "q8", Text: "Roughly what percentage of the Sun's mass is hydrogen?", Answer: 73},
	{ID: "q9", Text: "What percentage of the Moon's surface can be seen from Earth over time?", Answer: 59},
	{ID: "q10", Text: "The Moon's gravity is what percentage of Earth's?", Answer: 17},
	{ID: "q11", Text: "Mars' gravity is what percentage of Earth's?", Answer: 38},
	{ID: "q12", Text: "Roughly what percentage of the universe is ordinary matter (stars, planets, us)?", Answer: 5},
	{ID: "q13", Text: "Roughly what percentage of the universe is dark energy?", Answer: 68},

	// earth
	{ID: "q14", Text: "Roughly what percentage of Earth's atmosphere is argon?", Answer: 1},
	{ID: "q15", Text: "What percentage of Earth's water is in the oceans?", Answer: 97},
	{ID: "q16", Text: "Roughly what percentage of Earth's fresh water is locked in ice and glaciers?", Answer: 69},
	{ID: "q17", Text: "Roughly what percentage of the world's ice is in Antarctica?", Answer: 90},
	{ID: "q18", Text: "Antarctica covers roughly what percentage of Earth's land?", Answer: 9},
	{ID: "q19", Text: "Russia covers roughly what percentage of Earth's land?", Answer: 11},
	{ID: "q20", Text: "The Pacific Ocean covers roughly what percentage of Earth's surface?", Answer: 32},
	{ID: "q21", Text: "Roughly what percentage of Earth's crust (by mass) is oxygen?", Answer: 46},
	{ID: "q22", Text: "Roughly what percentage of Earth's crust (by mass) is silicon?", Answer: 28},
	{ID: "q23", Text: "The mantle makes up roughly what percentage of Earth's volume?", Answer: 84},
	{ID: "q24", Text: "Roughly what percentage of Greenland is covered by ice?", Answer: 80},
	{ID: "q25", Text: "The Sahara covers roughly what percentage of Africa?", Answer: 31},
	{ID: "q26", Text: "Roughly what percentage of the Amazon rainforest is in Brazil?", Answer: 60},
	{ID: "q27", Text: "Roughly what percentage of the Netherlands lies below sea level?", Answer: 26},
	{ID: "q28", Text: "Alaska is roughly what percentage of the total area of the United States?", Answer: 18},

	// people
	{ID: "q29", Text: "Roughly what percentage of the world's people live in Asia?", Answer: 59},
	{ID: "q30", Text: "Roughly what percentage of the world's people live in India?", Answer: 18},
	{ID: "q31", Text: "Roughly what percentage of Egyptians live along the Nile?", Answer: 95},
	{ID: "q32", Text: "Roughly what percentage of people are left-handed?", Answer: 10},
	{ID: "q33", Text: "Roughly what percentage of people in the US have type O blood?", Answer: 45},

	// body
	{ID: "q34", Text: "Roughly what percentage of DNA do humans share with chimpanzees?", Answer: 99},
	{ID: "q35", Text: "The brain is roughly what percentage of an adult's body weight?", Answer: 2},
	{ID: "q36", Text: "Roughly what percentage of the body's energy does the brain use?", Answer: 20},
	{ID: "q37", Text: "Roughly what percentage of the human body's mass is oxygen?", Answer: 65},
	{ID: "q38", Text: "What percentage of the body's calcium is in bones and teeth?", Answer: 99},
	{ID: "q39", Text: "What percentage of an adult's 206 bones are in the hands and feet?", Answer: 51},
	{ID: "q40", Text: "The skeleton is roughly what percentage of an adult's body weight?", Answer: 15},
	{ID: "q41", Text: "Roughly what percentage of the air you breathe out is oxygen?", Answer: 16},

	// food
	{ID: "q42", Text: "Roughly what percentage of a watermelon is water?", Answer: 92},
	{ID: "q43", Text: "Roughly what percentage of honey is water?", Answer: 17},
}
