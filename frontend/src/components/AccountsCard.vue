<template>
  <v-card class="px-4">
    <v-card-title>Accounts</v-card-title>
    <v-card-subtitle v-if="this.datasetId">{{
      this.datasetId
    }}</v-card-subtitle>
    <v-data-table
      :headers="headers"
      :items="entitlements"
      :search="search"
      :loading="loading"
    >
      <template v-slot:loading>
        <v-row justify="center" align="center">
          <div class="text-center ma-12">
            <v-progress-circular
              v-if="loading"
              indeterminate
              color="primary"
            ></v-progress-circular>
          </div>
        </v-row>
      </template>
      <template v-slot:top>
        <v-toolbar flat color="white" class="mb-1">
          <v-text-field
            v-model="search"
            clearable
            flat
            hide-details
            prepend-inner-icon="search"
            label="Search"
          ></v-text-field>
          <v-divider
            v-if="showAddAccountButton"
            class="mx-4"
            inset
            vertical
          ></v-divider>
          <v-btn
            v-if="showAddAccountButton"
            color="primary"
            dark
            @click.stop="addAccount"
            >Add Account</v-btn
          >
        </v-toolbar>
      </template>
      <template v-for="h in headers" v-slot:[`header.${h.value}`]="{ header }">
        <v-tooltip bottom v-bind:key="h.value">
          <template v-slot:activator="{ on }">
            <span v-on="on">{{ h.text }}</span>
          </template>
          <span v-if="header.tooltip">{{ header.tooltip }}</span>
        </v-tooltip>
      </template>
      <template v-slot:item.createdAt="{ item }">
        {{ toLocalTime(item.createdAt) }}
      </template>
      <template v-slot:item.policies="{ item }">
        <v-chip-group max="0" multiple column active-class="primary--text">
          <v-chip
            v-for="policy in item.policies.slice(0, 4)"
            :key="policy.policyId"
          >
            {{ policy.name }}
          </v-chip>
          <v-chip
            v-if="item.policies.length > 4 && showAddAccountButton"
            @click="editItem(item)"
          >
            ...
          </v-chip>
          <v-chip v-else-if="item.policies.length > 4">
            ...
          </v-chip>
        </v-chip-group>
      </template>
      <template v-if="this.showAddAccountButton" v-slot:item.action="{ item }">
        <v-tooltip top>
          <template v-slot:activator="{ on }">
            <v-icon v-on="on" small class="mr-2" @click="editItem(item)">
              edit
            </v-icon>
          </template>
          <span>Edit</span>
        </v-tooltip>
        <v-tooltip top>
          <template v-slot:activator="{ on }">
            <v-icon v-on="on" small class="mr-2" @click="deleteItem(item)">
              delete
            </v-icon>
          </template>
          <span>Delete</span>
        </v-tooltip>
        <v-tooltip top>
          <template v-slot:activator="{ on }">
            <v-icon
              v-on="on"
              v-if="item.marketplace.length > 0"
              small
              @click="resetItem(item)"
            >
              {{ icons.reset }}
            </v-icon>
          </template>
          <span>Reset Marketplace Account</span>
        </v-tooltip>
      </template>
    </v-data-table>
    <v-dialog v-model="showCreateAccount" persistent max-width="700px">
      <EditAccount
        :key="componentKey"
        :userData="selectedItem"
        v-on:close="editCompleted(true)"
        v-on:cancel="editCompleted(false)"
      />
    </v-dialog>
    <Dialog
      v-if="showDeleteDialog"
      v-model="showDeleteDialog"
      :title="deleteDialogTitle"
      :text="deleteDialogText"
      v-on:confirmed="deleteAccount(selectedItem)"
      v-on:canceled="deleteCompleted"
      confirmButtonColor="red darken-1"
      confirmButtonText="Delete"
    />
    <Dialog
      v-if="showResetDialog"
      v-model="showResetDialog"
      :title="resetDialogTitle"
      :text="resetDialogText"
      v-on:confirmed="resetAccount(selectedItem)"
      v-on:canceled="resetCompleted"
      confirmButtonColor="red darken-1"
      confirmButtonText="Reset"
    />
    <v-card-text
      style="height: 100px; position: relative"
      v-if="showAddAccountButton"
    >
      <v-fab-transition>
        <v-btn
          color="pink"
          dark
          absolute
          bottom
          right
          fab
          @click="addAccount()"
        >
          <v-icon>{{ icons.plus }}</v-icon>
        </v-btn>
      </v-fab-transition>
    </v-card-text>
  </v-card>
</template>

<script>
import EditAccount from '@/components/EditAccount';
import Dialog from '@/components/Dialog.vue';
import {
  mdiChevronRight,
  mdiForward,
  mdiDotsVertical,
  mdiPlus,
  mdiReplay
} from '@mdi/js';

export default {
  components: {
    EditAccount,
    Dialog
  },
  props: {
    selectedDataset: Object
  },
  data() {
    return {
      showAddAccountButton: true,
      entitlements: [],
      componentKey: 0,
      loading: false,
      showCreateAccount: false,
      showDeleteDialog: false,
      selectedItem: null,
      selectedItems: [],
      icons: {
        chevronRight: mdiChevronRight,
        forward: mdiForward,
        verticalDots: mdiDotsVertical,
        plus: mdiPlus,
        reset: mdiReplay
      },
      search: '',
      showResetDialog: false
    };
  },
  computed: {
    datasetId() {
      if (this.selectedDataset && this.selectedDataset.datasetId) {
        return this.selectedDataset.datasetId;
      }
      return null;
    },
    headers() {
      let h = [
        { text: 'Email', value: 'email', tooltip: 'The account email address' },
        {
          text: 'Email Type',
          value: 'emailType',
          tooltip: 'The email type of the account'
        },
        {
          text: 'Policies',
          value: 'policies',
          tooltip: 'The policies that the account is entitled to'
        },
        {
          text: 'Policy Search String',
          value: 'policySearchString',
          align: ' d-none'
        }
      ];

      if (!this.selectedDataset) {
        h.push(
          {
            text: 'Modified At',
            value: 'createdAt',
            tooltip: 'The last modified time for the account'
          },
          {
            text: 'Modified By',
            value: 'createdBy',
            tooltip: 'The last user that modified the account'
          },
          { text: '', value: 'action', sortable: false }
        );
      }

      return h;
    },
    deleteDialogTitle() {
      return 'Delete accounts?';
    },
    deleteDialogText() {
      return `Please click 'Delete' to confirm that you want to delete the selected account.`;
    },
    resetDialogTitle() {
      return 'Reset marketplace account?';
    },
    resetDialogText() {
      return `Please click 'Reset' to confirm that you want to delete the marketplace account(s) for ${this.selectedItem.email}.`;
    }
  },
  created() {
    if (this.selectedDataset) {
      this.showAddAccountButton = false;
    } else if (this.$route.query.datasetId) {
      this.selectedDataset = { datasetId: this.$route.query.datasetId };
    }
    this.loadAccounts();
  },
  methods: {
    addAccount() {
      this.componentKey += 1;
      this.selectedItem = null;
      this.selectedItems = [];
      this.showCreateAccount = true;
    },
    editItem(item) {
      this.componentKey += 1;
      this.selectedItem = item;
      this.showCreateAccount = true;
    },
    deleteCompleted() {
      this.selectedItem = null;
      this.selectedItems = [];
      this.showDeleteDialog = false;
    },
    editCompleted(refresh) {
      this.selectedItem = null;
      this.selectedItems = [];
      this.showCreateAccount = false;
      if (refresh) {
        this.loadAccounts();
      }
    },
    deleteItem(item) {
      this.selectedItem = item;
      this.selectedItems = [item];
      this.showDeleteDialog = true;
    },
    deleteItems() {
      this.showDeleteDialog = true;
    },
    loadAccounts() {
      this.loading = true;
      this.$store
        .dispatch('getAccounts', {
          datasetId: this.datasetId
        })
        .then(response => {
          if (response.success) {
            let data = response.data;
            data.forEach(function(element) {
              if (element.policies && element.policies.length > 0) {
                element.policySearchString = element.policies
                  .map(e => e.name)
                  .join('');
              } else {
                element.policySearchString = null;
              }
            });
            this.entitlements = data;
          } else {
            this.entitlements = [];
          }
          this.loading = false;
        });
    },
    deleteAccount(item) {
      this.loading = true;
      this.$store
        .dispatch('deleteAccount', {
          rowId: item.rowId,
          accountId: item.accountId
        })
        .then(() => {
          this.loadAccounts();
          this.loading = false;
        });
    },
    resetItem(item) {
      this.selectedItem = item;
      this.showResetDialog = true;
    },
    resetCompleted() {
      this.showResetDialog = false;
      this.selectedItem = null;
    },
    resetAccount() {
      this.$store
        .dispatch('submitProcurementAccountReset', {
          projectId: this.$store.state.settings.projectId,
          accountId: this.selectedItem.accountId
        })
        .then(result => {
          this.loading = false;
          if (!result.success) {
            this.showError = true;
            this.errorDialogTitle = 'Error resetting accounts';
            this.errorDialogText = result.errors.join(', ');
          } else {
            this.showError = false;
            this.loadAccounts();
            this.resetCompleted();
          }
        })
        .catch(error => {
          console.error(`Error submitting approval change: ${error}`);
        });
    },
    toLocalTime(epoch) {
      let d = new Date(epoch);
      return d.toLocaleString();
    }
  }
};
</script>
