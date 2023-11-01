import { onMounted, onUnmounted, ref } from 'vue'

export function usePinnable( pinClass ) {
   const toolbarTop = ref(0)
   const toolbarHeight = ref(0)
   const toolbarWidth = ref(0)
   const toolbar = ref(null)
   const pinned = ref(false)

   const scrolled = (() => {
      if ( window.scrollY <= toolbarTop.value ) {
         if ( toolbar.value.classList.contains("sticky") ) {
            toolbar.value.classList.remove("sticky")
            let dts = document.getElementsByClassName("p-datatable-wrapper")
            if ( dts ) {
               dts[0].style.top = `0px`
            }
         }
         pinned.value = false
      } else {
         if ( toolbar.value.classList.contains("sticky") == false ) {
            let dts = document.getElementsByClassName("p-datatable-wrapper")
            if ( dts ) {
               dts[0].style.top = `${toolbarHeight.value}px`
            }
            toolbar.value.classList.add("sticky")
            toolbar.value.style.width = `${toolbarWidth.value}px`
            pinned.value = true
         }
      }
   })

   onMounted( () => {
      let tb = null
      let tbs = document.getElementsByClassName( pinClass )
      if ( tbs ) {
         tb = tbs[0]
      }
      if ( tb) {
         toolbar.value = tb
         toolbarHeight.value = tb.offsetHeight
         toolbarWidth.value = tb.offsetWidth
         toolbarTop.value = tb.getBoundingClientRect().top
         window.addEventListener("scroll", scrolled)
      }
   })
   onUnmounted(() => window.removeEventListener("scroll", scrolled))

   return {toolbar, pinned}
}