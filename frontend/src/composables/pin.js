import { onMounted, ref, watch } from 'vue'
import { useWindowScroll, useElementBounding } from '@vueuse/core'

export function usePinnable( pinClass ) {
   const pinnedY = ref(0)
   const toolbar = ref()
   const pinned = ref(false)
   const { y } = useWindowScroll()
   const bounds = ref()

   watch(y, (newY) => {
      if ( pinned.value == false ) {
         if ( bounds.value.top <= 0 ) {
            pinnedY.value = y.value+bounds.value.top
            let dts = document.getElementsByClassName("p-datatable-table-container")
            if ( dts ) {
               dts[0].style.top = `${bounds.value.height}px`
               dts[0].style.position = 'relative'
               let panel = dts[0].closest('.p-panel-content')
               if (panel) {
                  let h = panel.clientHeight
                  panel.style.height = `${h}px`
               }
            }
            toolbar.value.classList.add("sticky")
            toolbar.value.style.width = `${bounds.value.width}px`
            pinned.value = true
         }
      } else {
         if ( newY <=  pinnedY.value) {
            toolbar.value.classList.remove("sticky")
            toolbar.value.style.width = `auto`
            let dts = document.getElementsByClassName("p-datatable-table-container")
            if ( dts ) {
               dts[0].style.top = `0px`
               dts[0].style.position = 'static'
               let panel = dts[0].closest('.p-panel-content')
               if (panel) {
                  panel.style.height = `auto`
               }
            }
            pinned.value = false
         }
      }
   })

   onMounted( () => {
      let tbs = document.getElementsByClassName( pinClass )
      if ( tbs ) {
         toolbar.value = tbs[0]
         bounds.value = useElementBounding( toolbar )
      }
   })
   return {}
}