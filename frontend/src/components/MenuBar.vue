<template>
   <nav class="menu">
      <ul role="menubar"
         @keydown.right.prevent.stop="nextMenu" @keyup.left.prevent.stop="prevMenu"
         @keydown.space.prevent.stop @keydown.down.prevent.stop @keydown.up.prevent.stop
         @keyup.esc="toggleSubMenu()"
      >
         <li role="none">
            <router-link role="menuitem" id="dashboardmenu" to="/">
               <span class="menu-item">Dashboard</span>
            </router-link>
         </li>
         <li role="none">
            <router-link role="menuitem" id="jobsmenu" to="/jobs" tabindex="-1" >
               <span class="menu-item">Job Statuses</span>
            </router-link>
         </li>

         <li role="none"
            @click.stop="toggleSubMenu('digitizationmenu')" @keydown.enter="toggleSubMenu()"
            @keydown.space="toggleSubMenu()" @keydown.up.prevent.stop="prevSubMenu"
            @keydown.down.prevent.stop="nextSubMenu"
         >
            <span role="menuitem" id="digitizationmenu" tabindex="-1"
               aria-haspopup="true" :aria-expanded="isOpen('digitizationmenu').toString()"
            >
               <span class="menu-item">Digitization</span>
               <i class="fas fa-caret-down submenu-arrow" :class="{ rotated: isOpen('digitizationmenu')}"></i>
            </span>
            <ul v-if="isOpen('digitizationmenu')" class="dropdown-menu"  id="digitization-menu"
               role="menu" @keydown.space.prevent.stop @keydown.enter.stop="linkClicked"
            >
               <li class="submenu" role="none">
                  <a :href="systemStore.projectsURL" role="menuitem" tabindex="-1" id="projectsub" target="_blank">Projects</a>
               </li>
               <li class="submenu" role="none">
                  <a :href="systemStore.reportsURL" role="menuitem" tabindex="-1" id="reportsub" target="_blank">Reports</a>
               </li>
            </ul>
         </li>

         <li role="none" class="right">
            <div role="menuitem" tabindex="-1"  id="signoutmenu" @click="signOut" @keydown.stop.enter="signOut">
               Sign out
            </div>
         </li>
      </ul>
   </nav>
</template>

<script setup>
import { useSystemStore } from "@/stores/system"
import { useUserStore } from "@/stores/user"
import { useRouter } from 'vue-router'
import { ref, onMounted, onUnmounted, nextTick } from 'vue'

const userStore = useUserStore()
const systemStore = useSystemStore()
const router = useRouter()

const noFocus = ref(false)
const menuBarIdx = ref(0)
const menuBar = ref([
   {id: "dashboardmenu", submenu:[], expanded: false},
   {id: "jobsmenu", submenu:[], expanded: false},
   {id: "digitizationmenu",
      submenu:["projectsub", "reportsub"],
      expanded: false, subMenuIdx: 0},
   {id: "signoutmenu", submenu:[], expanded: false},
])

onMounted(() => {
   window.addEventListener("click", resetMenus)
}),
onUnmounted(() => {
   window.removeEventListener("click", resetMenus)
})


function resetMenus() {
   menuBarIdx.value = 0
   closeSubMenus()
}
function closeSubMenus() {
   menuBar.value.forEach( mb => {
      if (mb.submenu.length > 0) {
         mb.expanded = false
         mb.idx = 0
      }
   })
}

function linkClicked() {
   noFocus.value = true
}
function isOpen( menuID ) {
   let m = menuBar.value.find( mb => mb.id == menuID)
   return m.expanded
}
function nextMenu() {
   closeSubMenus()
   menuBarIdx.value++
   if (menuBarIdx.value == menuBar.value.length) {
      menuBarIdx.value = 0
   }
   setMenuFocus()
}
function prevMenu() {
   closeSubMenus()
   menuBarIdx.value--
   if (menuBarIdx.value < 0) {
      menuBarIdx.value = menuBar.value.length - 1
   }
   setMenuFocus()
}
function nextSubMenu() {
   let currMenu = menuBar.value[menuBarIdx.value]
   if ( currMenu.submenu.length == 0) {
      return
   }
   if ( currMenu.expanded ) {
      currMenu.subMenuIdx++
   } else {
      currMenu.expanded = true
   }
   if ( currMenu.subMenuIdx == currMenu.submenu.length) {
      currMenu.subMenuIdx = 0
   }
   setMenuFocus()
}
function prevSubMenu() {
   let currMenu = menuBar.value[menuBarIdx.value]
   if ( currMenu.submenu.length == 0) {
      return
   }
   if ( currMenu.expanded ) {
      currMenu.subMenuIdx--
   } else {
      currMenu.expanded = true
   }
   if ( currMenu.subMenuIdx < 0) {
      currMenu.subMenuIdx = currMenu.submenu.length-1
   }
   setMenuFocus()
}
function toggleSubMenu( targetMenu ) {
   if (!targetMenu) {
       let menu = menuBar.value[menuBarIdx.value]
      if ( menu.submenu.length == 0) {
         return
      }
      menu.expanded = !menu.expanded
      menu.subMenuIdx = 0
      setMenuFocus()
   } else {
      menuBar.value.forEach( (m,idx) => {
         if (m.id != targetMenu) {
            m.expanded = false
         } else {
            m.expanded = !m.expanded
            m.subMenuIdx = 0
            menuBarIdx.value = idx
            setMenuFocus()
         }
      })
   }
}
function setMenuFocus() {
   if ( noFocus.value === true) {
      noFocus.value = false
      return
   }
   let menu = menuBar.value[menuBarIdx.value]
   if (menu.submenu.length == 0 || menu.expanded == false) {
      document.getElementById(menu.id).focus({preventScroll:true})
   } else {
      nextTick( () => {
         let eleID = menu.submenu[menu.subMenuIdx]
         let subMenuEle = document.getElementById(eleID)
         subMenuEle.focus({preventScroll:true})
      })
   }
}

function signOut() {
   userStore.signout()
   router.push("signedout")
}
</script>

<style scoped lang="scss">
.menu {
   background-color: var(--uvalib-blue-alt-darkest);
   padding: 5px 10px 0 10px;

   ul.dropdown-menu {
      position: absolute;
      z-index: 1000;
      background: white;
      padding: 0;
      border-radius: 0 0 5px 5px;
      border: 1px solid var(--uvalib-grey-light);
      border-top: none;
      overflow: hidden;
      transition: 200ms ease-out;
      display: grid;
      grid-auto-rows: auto;
      width: max-content;
      top: 29px;

      li.submenu {
         padding: 0;
         margin: 0;
         a, div {
            margin:0;
            padding: 10px 15px;
            font-weight: normal;
            color: var(--uvalib-text-dark);
            text-align: right;
            box-sizing: border-box;
            display: block;
            cursor: pointer;
            outline: none;
            top: 34px;
            &:hover {
               text-decoration: underline;
               background: var(--uvalib-blue-alt-light);
            }
            &:focus {
               background: var(--uvalib-blue-alt-light);
            }
         }
      }
   }

   ul {
      display: flex;
      position: relative;
      list-style: none;
      margin: 0;
      padding: 0;
      flex-flow: row wrap;
      justify-content: flex-start;
      li.right {
         margin-left: auto;
         padding-right: 0;
         cursor: pointer;
         &:hover {
            text-decoration: underline;
         }
      }
      li {
         display: inline-block;
         padding: 5px 15px 5px 0px;
         margin: 0;
         font-weight: 500;
         position: relative;
         color: white;
         a {
            color: white;
            text-decoration: none;
            &:hover {
               text-decoration: underline;
            }
         }
         .menu-item {
            cursor: pointer;
            color: white;
            flex: 0 1 auto;
            display: inline-block;
            border-bottom:1px solid transparent;
            &:hover {
               border-bottom:1px solid white;
            }
         }
         a.router-link-active {
            .menu-item{
               font-weight: bold;
               cursor:default;
               color: var(--uvalib-blue-alt-light);
               &:hover {
                  text-decoration: none;
                  border-bottom: none;
               }
            }
         }
      }
   }
   .submenu-arrow {
      transform: rotate(0deg);
      transition-duration: 200ms;
      display:inline-block;
      margin-left: 5px;
   }
     .submenu-arrow.rotated {
      transform: rotate(180deg);
   }
}
</style>