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
- **iOS Application**: https://github.com/RichSvK/StockBalance
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
- `GET /links` - Retrieve list of capital market-related websites link

### Search Stock
- `GET /stock` - Search stock with the given prefix text

## IPO Stock
- `GET /ipo` - Get IPO Stocks data
- `POST /ipo/condition` - Get IPO Stocks based on dynamic filter

## Broker
- `GET /brokers` - Get list of broker

## Stock Ownership
- `GET /balance/:code` - Get the stock ownership data
- `GET /balance/export` - Export stock ownership of the selected stock
- `POST /balance/upload` - Insert stock ownership KSEI data to the database
- `GET /api/auth/balance/scriptless` - Get a list of stocks with scriptless share changes over the past month
- `GET /api/auth/balance/change` - Get a list of stocks with ownership changes by shareholder type over the past month
