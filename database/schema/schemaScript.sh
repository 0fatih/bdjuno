echo "Creating cosmos schema"
psql -U bdjuno -d bdjuno -f 00-cosmos.sql

echo "Creating auth schema"
psql -U bdjuno -d bdjuno -f 01-auth.sql

echo "Creating bank schema"
psql -U bdjuno -d bdjuno -f 02-bank.sql

echo "Creating staking schema"
psql -U bdjuno -d bdjuno -f 03-staking.sql

echo "Creating consensus schema"
psql -U bdjuno -d bdjuno -f 04-consensus.sql

echo "Creating mint schema"
psql -U bdjuno -d bdjuno -f 05-mint.sql

echo "Creating distribution schema"
psql -U bdjuno -d bdjuno -f 06-distribution.sql

echo "Creating pricefeed schema"
psql -U bdjuno -d bdjuno -f 07-pricefeed.sql

echo "Creating gov schema"
psql -U bdjuno -d bdjuno -f 08-gov.sql

echo "Creating modules schema"
psql -U bdjuno -d bdjuno -f 09-modules.sql

echo "Creating slashing schema"
psql -U bdjuno -d bdjuno -f 10-slashing.sql

echo "Creating feegrant schema"
psql -U bdjuno -d bdjuno -f 11-feegrant.sql

echo "Creating upgrade schema"
psql -U bdjuno -d bdjuno -f 12-upgrade.sql

echo "Creating wasm schema"
psql -U bdjuno -d bdjuno -f 13-wasm.sql
