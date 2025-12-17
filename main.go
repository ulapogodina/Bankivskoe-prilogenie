package main

import (
 "errors"
 "fmt"
 "strconv"
 "time"
)

// Кастомные ошибки
var (
 ErrInsufficientFunds   = errors.New("недостаточно средств на счете")
 ErrInvalidAmount       = errors.New("некорректная сумма (отрицательная или нулевая)")
 ErrAccountNotFound     = errors.New("счет не найден")
 ErrSameAccountTransfer = errors.New("попытка перевода на тот же счёт")
)

// Интерфейс для работы со счетом
type AccountService interface {
 Deposit(amount float64) error
 Withdraw(amount float64) error
 Transfer(to *Account, amount float64) error
 GetBalance() float64
 GetStatement() string
}

// Интерфейс для работы с хранилищем
type Storage interface {
 SaveAccount(account *Account) error
 LoadAccount(accountID string) (*Account, error)
 GetAllAccounts() ([]*Account, error)
}

// Тип транзакции
type TransactionType string

const (
 Deposit  TransactionType = "DEPOSIT"
 Withdraw TransactionType = "WITHDRAW"
 Transfer TransactionType = "TRANSFER"
)

// Структура транзакции
type Transaction struct {
 Type      TransactionType
 Amount    float64
 Timestamp time.Time
 Message   string
}

// Структура счета
type Account struct {
 ID           string
 Owner        string
 Balance      float64
 Transactions []Transaction
}

// Реализация методов AccountService для Account
func (a *Account) Deposit(amount float64) error {
 if amount <= 0 {
  return ErrInvalidAmount
 }

 a.Balance += amount
 a.Transactions = append(a.Transactions, Transaction{
  Type:      Deposit,
  Amount:    amount,
  Timestamp: time.Now(),
  Message:   fmt.Sprintf("Пополнение: +%.2f", amount),
 })
 return nil
}

func (a *Account) Withdraw(amount float64) error {
 if amount <= 0 {
  return ErrInvalidAmount
 }
 if amount > a.Balance {
  return ErrInsufficientFunds
 }

 a.Balance -= amount
 a.Transactions = append(a.Transactions, Transaction{
  Type:      Withdraw,
  Amount:    amount,
  Timestamp: time.Now(),
  Message:   fmt.Sprintf("Снятие: -%.2f", amount),
 })
 return nil
}

func (a *Account) Transfer(to *Account, amount float64) error {
 if amount <= 0 {
  return ErrInvalidAmount
 }
 if a.ID == to.ID {
  return ErrSameAccountTransfer
 }
 if amount > a.Balance {
  return ErrInsufficientFunds
 }

 a.Balance -= amount
 to.Balance += amount

 // Запись транзакции для отправителя
 a.Transactions = append(a.Transactions, Transaction{
  Type:      Transfer,
  Amount:    amount,
  Timestamp: time.Now(),
  Message:   fmt.Sprintf("Перевод на счет %s: -%.2f", to.ID, amount),
 })

 // Запись транзакции для получателя
 to.Transactions = append(to.Transactions, Transaction{
  Type:      Transfer,
  Amount:    amount,
  Timestamp: time.Now(),
  Message:   fmt.Sprintf("Перевод от счета %s: +%.2f", a.ID, amount),
 })

 return nil
}

func (a *Account) GetBalance() float64 {
 return a.Balance
}

func (a *Account) GetStatement() string {
 if len(a.Transactions) == 0 {
  return "История транзакций пуста"
 }

 statement := "Выписка по счету:\n"
 statement += fmt.Sprintf("Владелец: %s\n", a.Owner)
 statement += fmt.Sprintf("Номер счета: %s\n", a.ID)
 statement += fmt.Sprintf("Текущий баланс: %.2f\n\n", a.Balance)
 statement += "История транзакций:\n"

 for i, tx := range a.Transactions {
  statement += fmt.Sprintf("%d. %s [%s]\n", i+1, tx.Message, tx.Timestamp.Format("2006-01-02 15:04:05"))
 }

 return statement
}

// InMemoryStorage реализация хранилища в памяти
type InMemoryStorage struct {
 accounts map[string]*Account
 nextID   int
}

func NewInMemoryStorage() *InMemoryStorage {
 return &InMemoryStorage{
  accounts: make(map[string]*Account),
  nextID:   1,
 }
}

func (s *InMemoryStorage) SaveAccount(account *Account) error {
 if account.ID == "" {
  account.ID = strconv.Itoa(s.nextID)
  s.nextID++
 }
 s.accounts[account.ID] = account
 return nil
}

func (s *InMemoryStorage) LoadAccount(accountID string) (*Account, error) {
 account, exists := s.accounts[accountID]
 if !exists {
  return nil, ErrAccountNotFound
 }
 return account, nil
}

func (s *InMemoryStorage) GetAllAccounts() ([]*Account, error) {
 accounts := make([]*Account, 0, len(s.accounts))
 for _, account := range
  nge s.accounts {
  accounts = append(accounts, account)
 }
 return accounts, nil
}

// BankApp - основное приложение
type BankApp struct {
 storage Storage
}

func NewBankApp(storage Storage) *BankApp {
 return &BankApp{storage: storage}
}

func (app *BankApp) CreateAccount(owner string) (*Account, error) {
 account := &Account{
  Owner:        owner,
  Balance:      0,
  Transactions: []Transaction{},
 }
 err := app.storage.SaveAccount(account)
 if err != nil {
  return nil, err
 }
 return account, nil
}

func (app *BankApp) FindAccount(accountID string) (*Account, error) {
 return app.storage.LoadAccount(accountID)
}

func (app *BankApp) GetAllAccounts() ([]*Account, error) {
 return app.storage.GetAllAccounts()
}

func main() {
 storage := NewInMemoryStorage()
 bankApp := NewBankApp(storage)

 for {
  fmt.Println("\n=== Банковское приложение ===")
  fmt.Println("1. Создать счет")
  fmt.Println("2. Пополнить счет")
  fmt.Println("3. Снять средства")
  fmt.Println("4. Перевести другому счету")
  fmt.Println("5. Просмотреть баланс")
  fmt.Println("6. Получить выписку")
  fmt.Println("7. Список всех счетов")
  fmt.Println("8. Выйти")
  fmt.Print("Выберите действие: ")

  var choice int
  fmt.Scan(&choice)

  switch choice {
  case 1:
   appCreateAccount(bankApp)
  case 2:
   appDeposit(bankApp)
  case 3:
   appWithdraw(bankApp)
  case 4:
   appTransfer(bankApp)
  case 5:
   appGetBalance(bankApp)
  case 6:
   appGetStatement(bankApp)
  case 7:
   appListAccounts(bankApp)
  case 8:
   fmt.Println("До свидания!")
   return
  default:
   fmt.Println("Неверный выбор")
  }
 }
}

func appCreateAccount(bankApp *BankApp) {
 var owner string
 fmt.Print("Введите имя владельца счета: ")
 fmt.Scan(&owner)

 account, err := bankApp.CreateAccount(owner)
 if err != nil {
  fmt.Printf("Ошибка создания счета: %v\n", err)
  return
 }

 fmt.Printf("Счет успешно создан!\n")
 fmt.Printf("Номер счета: %s\n", account.ID)
 fmt.Printf("Владелец: %s\n", account.Owner)
}

func appDeposit(bankApp *BankApp) {
 account, err := getAccount(bankApp)
 if err != nil {
  return
 }

 var amount float64
 fmt.Print("Введите сумму для пополнения: ")
 fmt.Scan(&amount)

 err = account.Deposit(amount)
 if err != nil {
  fmt.Printf("Ошибка пополнения: %v\n", err)
  return
 }

 bankApp.storage.SaveAccount(account)
 fmt.Printf("Счет успешно пополнен на %.2f. Новый баланс: %.2f\n", amount, account.GetBalance())
}

func appWithdraw(bankApp *BankApp) {
 account, err := getAccount(bankApp)
 if err != nil {
  return
 }

 var amount float64
 fmt.Print("Введите сумму для снятия: ")
 fmt.Scan(&amount)

 err = account.Withdraw(amount)
 if err != nil {
  fmt.Printf("Ошибка снятия: %v\n", err)
  return
 }

 bankApp.storage.SaveAccount(account)
 fmt.Printf("Со счета снято %.2f. Новый баланс: %.2f\n", amount, account.GetBalance())
}

func appTransfer(bankApp *BankApp) {
 fromAccount, err := getAccount(bankApp)
 if err != nil {
  return
 }

 var toAccountID string
 fmt.Print("Введите номер счета получателя: ")
 fmt.Scan(&toAccountID)

 toAccount, err := bankApp.FindAccount(toAccountID)
 if err != nil {
  fmt.Printf("Ошибка поиска счета получателя: %v\n", err)
  return
 }

 var amount float64
 fmt.Print("Введите сумму для перевода: ")
 fmt.Scan(&amount)

 err = fromAccount.Transfer(toAccount, amount)
 if err != nil {
  fmt.Printf("Ошибка перевода: %v\n", err)
  return
 }

 bankApp.storage.SaveAccount(fromAccount)
 bankApp.storage.SaveAccount(toAccount)
 fmt.Printf("Перевод успешно выполнен! Переведено %.2f на счет %s\n", amount, toAccountID)
}

func appGetBalance(bankApp *BankApp) {
 account, err := getAccount(bankApp)
 if err != nil {
  return
 }

 fmt.Printf("Текущий баланс: %.2f\n", account.GetBalance())
}

func appGetStatement(bankApp *BankApp) {
 account, err := getAccount(bankApp)
 if err != nil {
  return
 }

 fmt.Println(account.GetStatement())
}

func appListAccounts(bankApp *BankApp) {
 accounts, err := bankApp.GetAllAccounts()
 if err != nil {
  fmt.Printf("Ошибка получения списка счетов: %v\n", err)
  return
 }

 if len(accounts) == 0 {
  fmt.Println("Нет созданных счетов")
  return
 }

 fmt.Println("Список всех счетов:")
 for _, account := range accounts {
  fmt.Printf("Счет: %s, Владелец: %s, Баланс: %.2f\n", 
   account.ID, account.Owner, account.Balance)
 }
}

func getAccount(bankApp *BankApp) (*Account, error) {
 var accountID string
 fmt.Print("Введите номер счета: ")
 fmt.Scan(&accountID)

 account, err := bankApp.FindAccount(accountID)
 if err != nil {
  fmt.Printf("Ошибка поиска счета: %v\n", err)
  return nil, err
 }

 return account, nil
}
