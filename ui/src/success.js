export default {
  /* eslint-disable no-param-reassign */
  install(Vue) {
    const success = {
      deviceRename: 'renaming device',
      deviceDelete: 'deleting device',
      sessionClose: 'closing session',
      sessionRemoveRecord: 'deleting recorded session',
      firewallRuleCreating: 'creating rule',
      firewallRuleEditing: 'editing rule',
      firewallRuleDeleting: 'deleting rule',
      publicKeyCreating: 'creating public key',
      publicKeyEditing: 'editing public key',
      publicKeyDeleting: 'deleting public key',
      privateKeyCreating: 'creating private key',
      privateKeyEditing: 'editing private key',
      privateKeyDeleting: 'deleting private key',
      profileData: 'updating data',
      profilePassword: 'updating password',
      namespaceCreating: 'creating namespace',
      namespaceNewMember: 'adding new member',
      namespaceDelete: 'deleting namespace',
      namespaceEdit: 'editing namespace',
      namespaceRemoveUser: 'removing member',
      namespaceReload: 'reloading namespace',
      addUser: 'creating account',
      tokenList: 'list token',
      tokenCreating: 'creating token',
      tokenEditing: 'editing token',
      tokenDeleting: 'deleting token',
    };

    Vue.success = success;
    Vue.prototype.$success = success;
  },
  /* eslint-enable no-param-reassign */
};
