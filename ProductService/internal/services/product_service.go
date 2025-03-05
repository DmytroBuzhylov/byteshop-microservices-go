package services

var categoryList = []string{
	"books", "courses", "gaming-accounts", "digital-goods", "steam-keys", "online-services",
}

func CheckCategory(category string) bool {
	flag := false
	for _, i := range categoryList {
		if i == category {
			flag = true
		}
	}
	return flag
}
