<template>
  <div id="app">
    <b-container class="my-5">
      <template v-if="isAuthenticated">
        <div class="w-100 text-center mb-3">
          <b-form-select style="width: 20%" v-model="$i18n.locale">
            <b-form-select-option value="tr">TR</b-form-select-option>
            <b-form-select-option value="en">EN</b-form-select-option>
          </b-form-select>
        </div>
        <b-nav pills align="center">
          <b-nav-item exact exact-active-class="active" :to="{name: 'AlarmList'}">{{ $t("listAlarm") }}</b-nav-item>
          <b-nav-item exact exact-active-class="active" :to="{name: 'AlarmCreate'}">{{ $t("createAlarm") }}</b-nav-item>
        </b-nav>
        <hr>
        <router-view/>
      </template>
      <template v-else>
        <h5>{{ $t('telegramBotMsg') }}</h5>
        <a target="_blank" href="https://t.me/reminder_1996_bot">Bot Link</a>

        <div class="mt-3"></div>

        <b-form-group label="Token" label-for="token">
          <b-form-input id="token" v-model="token"></b-form-input>
        </b-form-group>

        <b-button @click="addToken" variant="outline-primary">OK</b-button>
      </template>
    </b-container>
  </div>
</template>

<script>
import Cookies from 'js-cookie'

export default {
  name: 'app',
  data() {
    return {
      lang: 'tr',
      isAuthenticated: false,
      token: '',
    }
  },
  methods: {
    async addToken() {
      if (this.token === "") {
        alert('Invalid token')
        return
      }

      try {
        const res = await fetch("/api/validate-token?token=" + this.token, {
          method: "GET",
        })

        const {code, message} = await res.json()
        if (code === 200) {
          this.isAuthenticated = true
          Cookies.set("token", this.token, {expires: 365})
          return
        }

        this.showErrorMessage(message)
      } catch (e) {
        console.error(e)
        alert(e.message)
      }
    }
  },
  created() {
    if (Cookies.get("token") && Cookies.get("token") !== "")
      this.isAuthenticated = true
  }
}
</script>

<style lang="scss">
.modal-body {
  white-space: pre-line !important;
}
</style>
