import { definePreset } from '@primevue/themes'
import Aura from '@primevue/themes/aura'
import ripple from '@primevue/themes/aura/ripple'
import tooltip from '@primevue/themes/aura/tooltip'
import './uva-colors.css'
import './styleoverrides.scss'

const UVA = definePreset(Aura, {
   root: {
      borderRadius: {
         none: '0',
         xs: '2px',
         sm: '3px',
         md: '4px',
         lg: '4px',
         xl: '8px'
      },
   },
   semantic: {
      primary: {
         50: 'var(--uvalib-brand-blue-lightest)',
         100: 'var(--uvalib-brand-blue-lighter)',
         200: 'var(--uvalib-brand-blue-lighter)',
         300: 'var(--uvalib-brand-blue-lighter)',
         400: 'var(--uvalib-brand-blue-light)',
         500: 'var(--uvalib-brand-blue-light)',
         600: 'var(--uvalib-brand-blue-light)',
         700: 'var(--uvalib-brand-blue-light)',
         800: 'var(--uvalib-brand-blue)',
         900: 'var(--uvalib-brand-blue)',
         950: 'var(--uvalib-brand-blue)'
      },
      focusRing: {
         width: '1px',
         style: 'dashed',
         offset: '2px'
      },
      formField: {
         paddingX: '0.75rem',
         paddingY: '0.5rem',
         borderRadius: '4px',
         focusRing: {
             width: '1px',
             style: 'dashed',
             color: '{primary.color}',
             offset: '2px',
             shadow: 'none'
         },
         transitionDuration: '{transition.duration}'
      },
      disabledOpacity: '0.3',
      colorScheme: {
         light: {
            primary: {
               color: '{primary.500}',
               contrastColor: '#ffffff',
               hoverColor: '{primary.100}',
               activeColor: '{primary.500}'
            },
            highlight: {
               background: '#ffffff',
               focusBackground: '#ffffff',
               color: 'var(--uvalib-text)',
               focusColor: '#ffffff'
            }
         },
      }
   },
   components: {
      button: {
         root: {
            paddingY: '.5em',
            paddingX: '1em',
         },
         colorScheme: {
            light: {
               secondary: {
                  background: 'var(--uvalib-grey-lightest)',
                  hoverBackground: 'var(--uvalib-grey-light)',
                  hoverBorderColor: 'var(--uvalib-grey)',
                  borderColor: 'var(--uvalib-grey-light)',
                  color: 'var(--uvalib-text)',
               },
            }
         }
      },
      datatable: {
         paginatorTop: {
            borderColor: 'var(--uvalib-grey-lightest)',
         },
         headerCell: {
            borderColor: 'var(--uvalib-grey-lightest)',
         },
         bodyCell: {
            borderColor: 'var(--uvalib-grey-lightest)',
        },
         colorScheme: {
            light: {
               root: {
                  borderColor: 'transparent'
              },
            }
         }
      },
      dialog: {
         colorScheme: {
            light: {
               root: {
                  background: '#ffffff',
                  borderColor: 'var(--uvalib-grey)',
                  padding: '15px',
                  borderRadius: '4px',
               },
               header: {
                  padding: '10px',
               },
               title: {
                  fontWeight: '600',
                  fontSize: '1em',
               }
            }
         }
      },
      menubar: {
         root: {
            background: "var(--uvalib-blue-alt-darkest)",
            borderColor: "var(--uvalib-blue-alt-darkest)",
            color: 'var(--uvalib-text)',
            borderRadius: "0px",
            padding: ".5rem .5rem",
            gap: '.25rem',
         },
         baseItem: {
            padding: "0.75rem 1rem"
         },
         item: {
            focusBackground: 'var(--uvalib-blue-alt)',
         },
         submenuIcon: {
            color: '#ffffff',
            focusColor: '#ffffff',
            activeColor: '{navigation.submenu.icon.active.color}'
        },
        submenu: {
            background: 'var(--uvalib-blue-alt-dark)',
            borderColor: 'var(--uvalib-blue-alt-dark)',
            borderRadius: '4px'
        }
      },
      panel: {
         header: {
            background: '#f8f9fa',
            borderColor:  'var(--uvalib-grey-light)',
            borderRadius: '4px 4px 0 0',
            padding: '1rem'
         },
         title: {
            fontWeight: '600',
         },
      },
      select: {
         root: {
            paddingY: '.5em',
            paddingX: '.5em',
            disabledBackground: '#fafafa',
            disabledColor: '#cacaca',
         },
         option: {
            selectedFocusBackground: 'var(--uvalib-blue-alt-light)',
            selectedFocusColor: 'var(--uvalib-text)',
            selectedBackground: 'var(--uvalib-blue-alt-light)',
            selectedColor: 'var(--uvalib-text)'
         }
      }
   },
   directives: {
      tooltip,
      ripple
   }
});

export default UVA;