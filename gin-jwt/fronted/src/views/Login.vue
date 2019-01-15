<template>
  <div class="login">
    <HelloWorld/>
  </div>
</template>

<script>
// @ is an alias to /src
import HelloWorld from '@/components/HelloWorld.vue';

export default {
  name: 'login',
  components: {
    HelloWorld,
  },
};
</script>
