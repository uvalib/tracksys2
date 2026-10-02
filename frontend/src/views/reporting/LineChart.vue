<template>
  <Line :options="options" :data="data" />
</template>

<script>
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
} from 'chart.js';
import { Line } from 'vue-chartjs'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
)

export default {
  name: 'LineChart',
  components: { Line },
  props: {
    data: {
        type: Object,
        required: true
      },
    options: {
      type: Object,
      default: () => {}
    }
  }
}
</script>