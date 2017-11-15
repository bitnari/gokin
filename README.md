# Gokin
A backend server for [Kakin](https://github.com/bitnari/kakin)


## API

**res 응답 코드**

코드 | 설명
--|--
0 | 성공
1 | 알 수 없는 오류
2 | 아이디의 포맷이 잘못됨
3 | 절대로 일어날 수 없는 오지고 지리고 아리랑 고개를 넘어 새가 날아들어 지저귀는 부분
4 | 비밀번호가 틀림
5 | 토큰의 포맷이 올바르지 않음
6 | 토큰이 만료되었거나 존재하지 않음
7 | 골드가 부족함

### 계정 정보
계정 인증을 할 때 사용합니다. 이 때 발급된 토큰은 10분 간 유효합니다.

```POST /api/v1/verify ```

**요청 파라미터**

키 | 설명 | 예시
--|--|--
grade | 학년 | 1
class | 반 | 1
id | 반 번호 | 1
password | 유저의 비밀번호 | specialpassword1234


**응답**

키 | 설명 | 예시
--|--|--
res | 성공 여부 | 0
token | 토큰 | (32자 랜덤 문자열)

**응답 예시**

```
Status 418 Teapot

{
	"res": 2
}
```


### 토큰 갱신
토큰을 더 길게 사용할 필요가 있을 경우에 사용합니다. 호출된 시간부터 10분 간 사용할 수 있게 됩니다.

```POST /api/v1/renew ```

**요청 파라미터**

키 | 설명 | 예시
--|--|--
token | 토큰 | (32자 랜덤 문자열)

**응답**

키 | 설명 | 예시
--|--|--
res | 성공 여부 | 0

### 계정 정보 얻기
계정의 돈을 얻을 때 사용합니다.

```POST /api/v1/account ```

**요청 파라미터**

키 | 설명 | 예시
--|--|--
token | 토큰 | (32자 랜덤 문자열)
password | 유저의 비밀번호 | specialpassword1234

**응답**

키 | 설명 | 예시
--|--|--
id | 유저의 아이디 | 10101
gold | 유저의 골드 | 0

**응답 예시**

```
Status 200 OK

{
	"id": "10101",
	"gold": 0
}
```


### 골드 사용
계정의 골드를 사용할 때 사용합니다.

```POST /api/v1/usegold ```

**요청 파라미터**

키 | 설명 | 예시
--|--|--
token | 토큰 | (32자 랜덤 문자열)
gold | 사용할 골드 | 0

**응답**

키 | 설명 | 예시
--|--|--
res | 성공 여부 | 0

**응답 예시**
```
Status 418 Teapot

{
	"res": 7
}
```