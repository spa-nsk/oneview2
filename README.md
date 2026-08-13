# oneview — Go-клиент HPE OneView REST API

Пакет для работы с:

- **HPE OneView appliance**, REST API **3800–8800** (OneView 6.60 и новее)
- **HPE OneView Global Dashboard**, REST API **2–300** по файлу `swagger 300.json`

Клиент сам определяет продукт по `GET /rest/version`, выставляет заголовок `X-API-Version`, логинится через `POST /rest/login-sessions` и шлёт токен в заголовке `Auth`.

Для appliance ответ логина — `sessionID`. Для Global Dashboard (swagger 300) — `token` + `user`. Поддерживаются оба варианта.

## Установка

```bash
go get github.com/spa-nsk/oneview2
```

## Быстрый старт

```go
package main

import (
    "context"
    "log"
    "github.com/spa-nsk/oneview2"
)

func main() {
    c, err := oneview.New(oneview.Config{
        Endpoint:    "https://oneview.example.com",
        Username:    "Administrator",
        Password:    "secret",
        Domain:      "Local",
        APIVersion:  0,    // 0 = currentVersion с appliance
        InsecureTLS: true, // только для лаборатории
    })
    if err != nil {
        log.Fatal(err)
    }
    ctx := context.Background()
    if err := c.Login(ctx); err != nil {
        log.Fatal(err)
    }
    defer c.Logout(ctx)

    hw, err := c.ListServerHardware(ctx, oneview.ListOptions{
        Count:  -1,
        Filter: []string{"status EQ 'OK'"},
        Sort:   "name:asc",
    })
    if err != nil {
        log.Fatal(err)
    }
    for _, s := range hw.Members {
        log.Println(s.Name, s.PowerState, s.Model)
    }
}
```

Переменные окружения (как у официального `oneview-golang`):

```
ONEVIEW_OV_ENDPOINT=https://oneview.example.com
ONEVIEW_OV_USER=Administrator
ONEVIEW_OV_PASSWORD=secret
ONEVIEW_OV_DOMAIN=Local
ONEVIEW_APIVERSION=3800
ONEVIEW_SSLVERIFY=false
```

```go
c, err := oneview.New(oneview.ConfigFromEnv())
```

Пример: `examples/list_inventory`.

## Выгрузка конфигурации серверов

Верхнеуровневые функции собирают identity, процессоры, DIMM, RAID/диски, PCI-устройства, NIC-порты, firmware, BIOS, профиль и enclosure.

```go
exp, err := c.ExportServer(ctx, "Encl1, bay 5", oneview.ExportOptions{})
if err != nil {
    log.Fatal(err)
}
fmt.Print(exp.Summary())
_ = oneview.SaveServerExportJSON("server.json", exp)

// все серверы
all, err := c.ExportServers(ctx, oneview.ListOptions{Count: -1}, oneview.ExportOptions{})
```

Идентификатор: имя, `serverName`, серийный номер, UUID или URI `/rest/server-hardware/...`.

Источники: `GET /rest/server-hardware/{id}` (`subResources`: Memory, LocalStorage, Devices), плюс `/firmware`, `/localStorageV2`, `/bios`, `/processors`, `/memory`.

Примеры:

```bash
go run ./examples/export_server "Encl1, bay 5"
go run ./examples/export_servers ./out
```

## Версии API

| Продукт | X-API-Version |
| --- | --- |
| Global Dashboard 2.20–3.00 (`swagger 300.json`) | 300 |
| OneView 6.60 | 3800 |
| OneView 7.00 | 4000 |
| OneView 8.00 | 4600 |
| OneView 8.20 | 5000 |
| OneView 9.00 | 6600 |
| OneView 10.00 | 7600 |
| OneView 11.00 / 11.10 / 11.20 | 8200 |
| OneView 11.30 | 8600 |
| (диапазон пакета) | **3800–8800** |

Алгоритм как в документации HPE:

1. `GET /rest/version`
2. проверить `minimumVersion <= requested <= currentVersion`
3. дальше все вызовы с `X-API-Version: <requested>`

Если `APIVersion` не задан, берётся `currentVersion` appliance.

JSON от Global Dashboard (API 300) и appliance (3800–8800) расходится по типам: `eTag` — строка или число, `powerLock` — bool или строка, `hostOsType` — int / `null` / имя ОС, `position` и счётчики CPU/RAM — int или `null` на rack-серверах. Клиент приводит это к `FlexString` / `FlexInt` / `FlexBool` и через `DecodeJSON` принимает оба варианта.

## Что покрыто

### Из swagger 300.json (Global Dashboard)

| Ресурс | Методы |
| --- | --- |
| `/rest/version` | GET |
| `/rest/login-sessions` | POST, DELETE |
| `/rest/resource-alerts` | GET list/id |
| `/rest/admin-settings/alert-settings` | GET, PUT |
| `/rest/appliances` | GET, POST, PATCH, DELETE, SSO |
| `/rest/audit-logs/settings` | GET, PUT, test-forwarding |
| `/rest/certificates/...` | GET remote, POST/DELETE servers |
| `/rest/converged-systems` | GET |
| `/rest/datacenters` | GET |
| `/rest/drive-enclosures` | GET |
| `/rest/enclosures` | GET, oaSsoUrl |
| `/rest/groups` | GET, POST, PATCH, DELETE, members |
| `/rest/interconnects` | GET |
| `/rest/logical-interconnects` | GET |
| `/rest/managed-sans` | GET |
| `/rest/appliance/network-interfaces` | GET, POST |
| `/rest/appliance/configuration/time-locale` | GET, POST |
| `/rest/san-managers` | GET |
| `/rest/server-firmware` | GET |
| `/rest/server-hardware` | GET, iloSsoUrl |
| `/rest/server-profiles` | GET |
| `/rest/server-profile-templates` | GET |
| `/rest/storage-pools` | GET |
| `/rest/storage-systems` | GET |
| `/rest/storage-volumes` | GET |
| `/rest/tasks` | GET, wait |

### OneView appliance 3800–8800 (сверх swagger)

CRUD и операции для профилей, ethernet/FC/FCoE сетей, network sets, enclosure groups, logical enclosures, LIG, uplink sets, firmware drivers, volumes (`/rest/volumes`), storage systems, power/refresh серверов, compliance LI.

Любой недокументированный путь:

```go
var raw map[string]any
err := c.GetJSON(ctx, "/rest/ethernet-networks?count=-1", &raw)
resp, err := c.PostJSON(ctx, "/rest/ethernet-networks", body, &created)
```

Асинхронные вызовы (HTTP 202 + `Location`):

```go
resp, task, err := c.CreateServerProfile(ctx, profile, true) // wait=true
```

## Фильтры

```go
oneview.ListOptions{
    Start:  0,
    Count:  50, // -1 = все (до лимита страницы)
    Filter: []string{"status EQ 'OK'", "name EQ 'profile1'"},
    Query:  "name EQ 'resource1'",
    Sort:   "name:asc",
    Fields: "name,uri,status",
}
```

На Global Dashboard можно ограничить список `GroupURIs`.
