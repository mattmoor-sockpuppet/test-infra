#!/bin/bash

# Copyright 2019 The Knative Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# GET=1 (attempt to fetch Postgres database and Grafana database from the test server)
# SKIPDBS=1 (entirely skip project's database operations)
# SKIPADDALL=1 (skip adding/merging to allprj)
# SKIPGRAFANA=1 (skip all grafana related stuff)
set -o pipefail
if [ -z "$PG_PASS" ]
then
  echo "$0: You need to set PG_PASS environment variable to run this script"
  exit 1
fi
if ( [ -z "$PROJ" ] || [ -z "$PROJDB" ] || [ -z "$PROJREPO" ] || [ -z "$PORT" ] || [ -z "$GA" ] || [ -z "$ICON" ] || [ -z "$ORGNAME" ] || [ -z "$GRAFSUFF" ] )
then
  echo "$0: You need to set PROJ, PROJDB, PROJREPO, PORT, GA, ICON, ORGNAME, GRAFSUFF environment variables to run this script"
  exit 2
fi
function finish {
    sync_unlock.sh
    rm -f /tmp/deploy.wip 2>/dev/null
}
if [ -z "$TRAP" ]
then
  sync_lock.sh || exit -1
  trap finish EXIT
  export TRAP=1
  > /tmp/deploy.wip
fi
echo "$0: $PROJ deploy started"
if [ -z "$SKIPDBS" ]
then
  PDB=1 TSDB=1 ./devel/create_databases.sh || exit 3
fi
if [ -z "$SKIPGRAFANA" ]
then
  ./devel/create_grafana.sh || exit 5
fi
echo "$0: $PROJ deploy finished"
