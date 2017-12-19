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