import { definePreset } from '@primeuix/themes'
import Aura from '@primeuix/themes/aura'
import ripple from '@primeuix/themes/aura/ripple'
import tooltip from '@primeuix/themes/aura/tooltip'
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
      accordion: {
         header: {
            background: '#f8f9fa',
            hoverBackground: '#f5f5ff',
            activeBackground: '#f8f9fa',
            activeHoverBackground: '#f8f9fa',
            borderRadius: 0,
         },
         panel: {
            borderWidth: '1px',
            borderColor: 'var(--uvalib-grey-lightest)',
            hoverBackground: 'var(--uvalib-grey-lightest)',
        },
        content: {
            background: '#ffffff',
            borderWidth: '1px 0 0 0',
            padding: '1.125rem 1.125rem 1.125rem 1.125rem'
        }
      },
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
         row: {
            selectedBackground: 'var(--uvalib-blue-alt-light)',
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
            borderRadius: "0px",
            background: 'var(--uvalib-blue-alt-darkest)',
            borderColor: 'var(--uvalib-blue-alt-darkest)',
         },
         baseItem: {
            borderRadius: '4px',
         },
         item: {
            focusBackground: 'var(--uvalib-blue-alt)',
            activeBackground: 'var(--uvalib-blue-alt)',
            color: '#ffffff',
            focusColor: '#ffffff',
            activeColor: '#ffffff',
            icon: {
               color: '#ffffff',
               focusColor: '#ffffff',
               activeColor: '#ffffff',
           }
         },
         submenu: {
            background: 'var(--uvalib-blue-alt-dark)',
            borderColor: 'var(--uvalib-blue-alt-dark)',
            borderRadius: '4px',
            shadow: '{overlay.navigation.shadow}',
            mobileIndent: '1rem',
        },
        submenuIcon: {
            size: '{navigation.submenu.icon.size}',
            color: '#ffffff',
            focusColor: '#ffffff',
            activeColor: '#ffffff',
        },
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
      },
      togglebutton: {
         colorScheme: {
            light: {
               root: {
                  background: '#ffffff',
                  checkedBackground: 'var(--uvalib-blue-alt-light)',
                  hoverBackground: '{surface.100}',
                  borderColor: 'var(--uvalib-grey-light)',
                  color: 'var(--uvalib-text)',
                  checkedColor: 'var(--uvalib-text)',
                  checkedBorderColor: 'var(--uvalib-blue-alt-light)',
               },
            }
         }
      },
   },
   directives: {
      tooltip,
      ripple
   }
});

export default UVA;