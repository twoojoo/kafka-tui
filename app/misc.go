package app

func GetTitle() string {
	title := " ╷  _  _______ _   _ ___"
	title += "\n │ | |/ /_   _| | | |_ _|"
	title += "\n │ | ' /  | | | | | || |"
	title += "\n │ | . \\  | | | |_| || |  v" + Version
	title += "\n │ |_|\\_\\ |_|  \\___/|___| by twoojoo"
	title += "\n └───────────────────────────────────────"
	return title
}

func GetHotkeysText1() string {
	htkTxt := "\nMove through tabs"
	// htkTxt += "\nFocus search bar"
	htkTxt += "\nSelect element"
	htkTxt += "\nMove"
	htkTxt += "\nAdd topic"

	return htkTxt
}

func GetHotkeysKeys1() string {
	htkTxt := "\nTab    "
	// htkTxt += "\n\\   "
	htkTxt += "\nEnter    "
	htkTxt += "\n🡡 🡣    "
	htkTxt += "\na    "

	return htkTxt
}

func GetHotkeysText2() string {
	htkTxt := "\nEdit topic config"
	htkTxt += "\nClear topic messages"
	htkTxt += "\nRecreate topic"
	htkTxt += "\nDelete topic"

	return htkTxt
}

func GetHotkeysKeys2() string {
	htkTxt := "\ne   "
	htkTxt += "\nc   "
	htkTxt += "\nr   "
	htkTxt += "\nd   "

	return htkTxt
}

func GetKafkaLogo() string {
	logo := "\n\n"
	logo += "                    @@@@@@\n"
	logo += "                   @@    @@@\n"
	logo += "                   @@    @@@\n"
	logo += "                    @@@@@@     @@@@@@@@\n"
	logo += "                      @@      @@@    @@@\n"
	logo += "                   @@@@@@@@  %@@@@  @@@,\n"
	logo += "                 #@@@    .@@@%   &@@@\n"
	logo += "                 @@@       @@&\n"
	logo += "                  @@@    @@@@@   @@@@\n"
	logo += "                   @@@@@@@   @@@%   @@@\n"
	logo += "                      @@      @@@    @@@\n"
	logo += "                    @@@@@@*    %@@@@@@@\n"
	logo += "                   @@    @@@\n"
	logo += "                   @@    @@@\n"
	logo += "                    @@@@@@"
	return logo
}
