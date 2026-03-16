## Stock Backend Service

### Program Description
Stock backend is a backend systems for stock ownership service and IPO stock service.

### Data Source
- Stock Ownership Data
  Source: KSEI (Kustodian Sentral Efek Indonesia)<br>
  URL: https://ksei.co.id/archive_download/holding_composition
- IPO Sample Data
  Source: E-IPO Stock Prospectus<br>
  URL: https://e-ipo.co.id/en

## Related Repositories
- **StockWeb Frontend**: https://github.com/RichSvK/StockWeb
- **iOS Application**: https://github.com/RichSvK/Stockbalances
- **API Gateway**: https://github.com/RichSvK/API_Gateway
- **User and Watchlist services**: https://github.com/RichSvK/User_Backend
- **GoIPO**: https://github.com/RichSvK/GoIPO

### System Requirements
Software used in developing this program:
- Go 1.24
- Gin Web Framework
- MySQL

## API Endpoints
### Website Link
- `GET /api/v1/links` - Retrieve list of capital market-related websites link
- `POST /api/v1/links` - Create a new website link notes
- `PATCH /api/v1/links` - Update an existing website link
- `DELETE /api/v1/links/:name` - Delete a website link by name

### Search Stock
- `GET /api/v1/stocks` - Search stock with the given prefix text

## IPO Stock
- `GET /api/v1/ipo` - Get IPO Stocks data
- `POST /api/v1/ipo/condition` - Get IPO Stocks based on dynamic filter

## Broker
- `GET /api/v1/brokers` - Get list of broker or filtered broker by condition

## Stock Ownership
- `GET /api/v1/balances/:code` - Get the stock ownership data
- `GET /api/v1/balances/export` - Export stock ownership of the selected stock
- `POST /api/v1/balances/import` - Insert stock ownership KSEI data to the database
- `GET /api/v1/auth/balances/scriptless` - Get a list of stocks with scriptless share changes over the past month
- `GET /api/auth/balances/changes` - Get a list of stocks with ownership changes by shareholder type over the past month