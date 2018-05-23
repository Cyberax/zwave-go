#include "Manager.h"

namespace OpenZWave {
	class CallbackAdapter {
	public:
		virtual ~CallbackAdapter() {};
		virtual void onManagerNotify(Notification const *notify) = 0;

		static void handleManagerNotification(Notification const* _pNotification, void* _context) {
			reinterpret_cast<CallbackAdapter*>(_context)->onManagerNotify(_pNotification);
		}
	};

	inline void addNotifyHandler(Manager *manager, CallbackAdapter *adapter) {
		manager->AddWatcher(&CallbackAdapter::handleManagerNotification, adapter);
	}
	
	inline void removeNotifyHandler(Manager *manager, CallbackAdapter *adapter) {
		manager->RemoveWatcher(&CallbackAdapter::handleManagerNotification, adapter);
	}
};
