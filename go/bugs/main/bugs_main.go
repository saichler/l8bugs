/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 *
 * Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
 * You may obtain a copy of the License at:
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"github.com/saichler/l8bugs/go/bugs/common"
	"github.com/saichler/l8bugs/go/bugs/services"
	"github.com/saichler/l8bus/go/overlay/vnic"
	l8common "github.com/saichler/l8common/go/common"
	"github.com/saichler/l8types/go/ifs"
	"os/exec"
	"time"
)

func main() {
	res := l8common.CreateResources("BugsServices", "/data/logs/l8bugs", uint32(common.BUGS_VNET))
	ifs.SetNetworkMode(ifs.NETWORK_K8s)
	nic := vnic.NewVirtualNetworkInterface(res, nil)
	nic.Start()
	nic.WaitForConnection()

	//Start postgres
	startDb(nic)

	services.ActivateBugsServices(common.DB_CREDS, common.DB_NAME, nic)

	l8common.WaitForSignal(res)
}

func startDb(nic ifs.IVNic) {
	_, user, pass, _, err := nic.Resources().Security().Credential(common.DB_CREDS, common.DB_NAME, nic.Resources())
	if err != nil {
		panic(common.DB_CREDS + " " + err.Error())
	}
	// When there is no security provider
	if user == "admin" && pass == "admin" {
		common.DB_NAME = "admin"
	}
	fmt.Println("/start-postgres.sh", common.DB_NAME, user, pass)
	var cmd *exec.Cmd
	cmd = exec.Command("nohup", "/start-postgres.sh", common.DB_NAME, user, pass)
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
	time.Sleep(time.Second * 5)
}
