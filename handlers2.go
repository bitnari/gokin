package main

import (
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

func Pay_v2(e echo.Context) error {
	token := e.FormValue("token")
	if len(token) != 32 {
		return e.JSON(http.StatusBadRequest, R {
			"res": ResErrInvalidToken,
		})
	}

	t, err := tokens.Get(token)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, R {
			"res": ResErrNoToken,
		})
	}

	account, err := mongo.GetAccount(t.User)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, R {
			"res": ResErrUnknown,
		})
	}

	credit, err := strconv.Atoi(e.FormValue("credit"))
	if err != nil {
		return e.JSON(http.StatusBadRequest, R {
			"res": ResErrInvalidCredit,
		})
	}

	if account.Gold + account.Credit < credit {
		return e.JSON(http.StatusTeapot, R {
			"res": ResErrNoGold,
		})
	}

	if config.Server.CreditPriority == "credit" {
		account.Credit -= credit
		if account.Credit < 0 {
			account.Gold += account.Credit
			account.Credit = 0
		}
	}else{
		account.Gold -= credit
		if account.Gold < 0 {
			account.Credit += account.Gold
			account.Gold = 0
		}
	}

	err = mongo.SetCredit(account.Id, account.Gold, account.Credit)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, R {
			"res": ResErrUnknown,
		})
	}

	return e.JSON(http.StatusOK, R {
		"res": ResSuccess,
	})
}

func VerifyAccount_v2(e echo.Context) error {
	id := toId(e.FormValue("grade"), e.FormValue("class"), e.FormValue("id"))
	if len(id) != 5 {
		return e.JSON(http.StatusBadRequest, R {
			"res": ResErrIdLenMismatch, // 유저이름의 문자열 길이가 틀리다
		})
	}
	password := e.FormValue("password")

	err := mongo.VerifyAccount(id, password)
	if err != nil {
		if err == ErrNoAccount {
			return e.JSON(http.StatusUnauthorized, R {
				"res": ResErrNoAccount,
			})
		}else if err == ErrIncorrectPassword {
			return e.JSON(http.StatusUnauthorized, R {
				"res": ResErrIncorrectPassword,
			})
		}else{
			return e.JSON(http.StatusInternalServerError, R {
				"res": ResErrUnknown,
			})
		}
	}

	return e.JSON(http.StatusOK, R {
		"res": ResSuccess, // 성공
		"token": tokens.New(id).Token,
	})
}

func Account_v2(e echo.Context) error {
	token := e.FormValue("token")
	if len(token) != 32 {
		return e.JSON(http.StatusBadRequest, R {
			"res": ResErrInvalidToken,
		})
	}

	t, err := tokens.Get(token)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, R {
			"res": ResErrNoToken,
		})
	}

	account, err := mongo.GetAccount(t.User)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, R {
			"res": ResErrUnknown,
		})
	}

	return e.JSON(http.StatusOK, R {
		"res": ResSuccess,
		"id": account.Id,
		"gold": account.Gold,
		"credit": account.Credit,
	})
}

func Score_v2(e echo.Context) error {
	token := e.FormValue("token")
	if len(token) != 32 {
		return e.JSON(http.StatusBadRequest, R {
			"res": ResErrInvalidToken,
		})
	}

	t, err := tokens.Get(token)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, R {
			"res": ResErrNoToken,
		})
	}

	account, err := mongo.GetAccount(t.User)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, R {
			"res": ResErrUnknown,
		})
	}

	score, err := strconv.Atoi(e.FormValue("score"))
	if err != nil {
		return e.JSON(http.StatusBadRequest, R {
			"res": ResErrInvalidCredit,
		})
	}

	gameId := e.FormValue("gameId")
	if !config.HasGame(gameId) {
		return e.JSON(http.StatusForbidden, R {
			"res": ResErrInvalidRequest,
		})
	}

	err = mongo.SetScore(account, gameId, score)

	return e.JSON(http.StatusOK, R {
		"res": ResSuccess,
	})
}
