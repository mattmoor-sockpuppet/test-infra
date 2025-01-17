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

if [ -z "$ONLY" ]
then
  all="knative"
else
  all=$ONLY
fi
mkdir sqlite 1>/dev/null 2>/dev/null
touch sqlite/touch
for proj in $all
do
    db=$proj
    rm -f sqlite/* 2>/dev/null
    touch sqlite/touch
    ./sqlitedb /var/lib/grafana.$db/grafana.db || exit 1
    rm -f grafana/dashboards/$proj/*.json || exit 2
    mv sqlite/*.json grafana/dashboards/$proj/ || exit 3
done
echo 'OK'
