import Vuex from 'vuex';
import { mount, createLocalVue } from '@vue/test-utils';
import DeviceRejectedList from '@/components/device/DeviceRejectedList';
import Vuetify from 'vuetify';

describe('DeviceRejectedList', () => {
  const localVue = createLocalVue();
  const vuetify = new Vuetify();
  localVue.use(Vuex);
  localVue.filter('moment', () => {});

  let wrapper;

  const numberDevices = 1;

  const pagination = {
    groupBy: [],
    groupDesc: [],
    itemsPerPage: 10,
    multiSort: false,
    mustSort: false,
    page: 1,
    sortBy: [],
    sortDesc: [],
  };

  const headers = [
    {
      text: 'Hostname',
      value: 'hostname',
      align: 'center',
    },
    {
      text: 'Operating System',
      value: 'info.pretty_name',
      align: 'center',
      sortable: false,
    },
    {
      text: 'Request Time',
      value: 'request_time',
      align: 'center',
      sortable: false,
    },
    {
      text: 'Actions',
      value: 'actions',
      align: 'center',
      sortable: false,
    },
  ];

  const devices = [
    {
      uid: '2378hj238',
      name: '37-23-hf-1c',
      identity: {
        mac: '00:00:00:00:00:00',
      },
      info: {
        id: 'linuxmint',
        pretty_name: 'Linux Mint 20.0',
        version: '',
      },
      public_key: '---pub_key---',
      tenant_id: '8490393000',
      last_seen: '2020-05-22T18:58:53.276Z',
      online: true,
      namespace: 'user',
      status: 'rejected',
    },
  ];

  const store = new Vuex.Store({
    namespaced: true,
    state: {
      devices,
      numberDevices,
    },
    getters: {
      'devices/list': (state) => state.devices,
      'devices/getNumberDevices': (state) => state.numberDevices,
    },
    actions: {
      'modals/showAddDevice': () => {
      },
      'devices/fetch': () => {
      },
      'devices/rename': () => {
      },
      'devices/resetListDevices': () => {
      },
      'stats/get': () => {
      },
    },
  });

  beforeEach(() => {
    wrapper = mount(DeviceRejectedList, {
      store,
      localVue,
      stubs: ['fragment', 'router-link'],
      vuetify,
    });
  });

  it('Is a Vue instance', () => {
    expect(wrapper).toBeTruthy();
  });
  it('Renders the component', () => {
    expect(wrapper.html()).toMatchSnapshot();
  });
  it('Compare data with default value', () => {
    expect(wrapper.vm.pagination).toEqual(pagination);
    expect(wrapper.vm.headers).toEqual(headers);
  });
  it('Process data in the computed', () => {
    expect(wrapper.vm.getListRejectedDevices).toEqual(devices);
    expect(wrapper.vm.getNumberRejectedDevices).toEqual(numberDevices);
  });
  it('Renders the template with components', async () => {
    expect(wrapper.find('[data-test="deviceIcon-component"]').exists()).toEqual(true);
    expect(wrapper.find('[data-test="DeviceActionButtonAccept-component"]').exists()).toEqual(true);
    expect(wrapper.find('[data-test="deviceActionButtonReject-component"]').exists()).toEqual(true);
  });
  it('Renders the template with data', () => {
    const dt = wrapper.find('[data-test="dataTable-field"]');
    const dataTableProps = dt.vm.$options.propsData;

    expect(dataTableProps.items).toHaveLength(numberDevices);
  });
});
