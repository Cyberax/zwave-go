// pair.i - SWIG interface
%module openzwave
     
%feature("flatnested");

%rename(inc) operator++;
%rename(eq) operator==;
%rename(neq) operator!=;

#define OPENZWAVE_EXPORT 
#define OPENZWAVE_EXPORT_WARNINGS_OFF
#define OPENZWAVE_EXPORT_WARNINGS_ON
#define DEPRECATED

namespace std
{
  %ignore runtime_error;
  struct runtime_error {};
}

%{
    #include "Bitfield.h"
    #include "Driver.h"
    #include "Manager.h"
    #include "Node.h"
    #include "Defs.h"
    #include "Options.h"
    #include "Utils.h"   
    #include "DoxygenMain.h"
    #include "Group.h"
    #include "Msg.h"
    #include "Notification.h"
    #include "OZWException.h"
    #include "Scene.h"
    #include "ZWSecurity.h"
%}


// Parse the original header file
%include "Bitfield.h"
%include "Driver.h"
%include "Manager.h"
%include "Node.h"
%include "Defs.h"
%include "Options.h"
%include "Utils.h"   
%include "DoxygenMain.h"
%include "Group.h"
%include "Msg.h"
%include "Notification.h"
%include "OZWException.h"
%include "Scene.h"
%include "ZWSecurity.h"
