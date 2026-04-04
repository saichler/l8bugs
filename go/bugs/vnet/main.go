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
	"github.com/saichler/l8bugs/go/bugs/common"
	l8common "github.com/saichler/l8common/go/common"
	"os"

	"github.com/saichler/l8bus/go/overlay/vnet"
)

func main() {
	resources := l8common.CreateResources("vnet-"+os.Getenv("HOSTNAME"), "/data/logs/l8bugs", uint32(common.BUGS_VNET))
	net := vnet.NewVNet(resources)
	net.Start()
	resources.Logger().Info("vnet started!")
	l8common.WaitForSignal(resources)
}
