# Bankivskoe-prilogenie
Bankivskoe-prilogenie
Требования по коду. Программа должна содержать 2 интерфейса:

AccountService - основной интерфейс для работы со счетом: type AccountService interface { Deposit (amount float64) error Withdraw (amount float64) error Transfer (to *Account, amount float64) error GetBalance () float64 GetStatement() string }
Storage - интерфейс для работы с хранилищем данных (инкапсулирует логику сохранения/загрузки): type Storage interface { SaveAccount (account *Account) error LoadAccount (accountID string) (*Account, error) GetAllAccounts () ([] *Account, error) } Кастомные ошибки. Создайте набор кастомных ошибок с помощью errors.New() или fmt.Errorf():  ErrInsufficientFunds - недостаточно средств на счете;  ErrInvalidAmount - некорректная сумма (отрицательная или нулевая);  ErrAccountNotFound - счет не найден;  ErrSameAccountTransfer - попытка перевода на тот же счёт.   Функциональные требования  создать счет (запрашивает имя владельца);  пополнить счет;  снять средства;  перевести другому счету (требует ввода ID счета);  просмотреть баланс;  получить выписку (показать историю транзакций);  выйти.
