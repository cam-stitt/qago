Step
====

A singular step to be run within a test case.


name
----

================= ==========
**required**      false
**type**          string
**default**       ""
================= ==========

The name of the step. This will be output for reference of the current step being run.

Example
^^^^^^^

.. literalinclude:: example.yml
   :language: yaml
   :emphasize-lines: 1
   :lines: 6-8


predicate
---------

================= =================
**required**      false
**type**          :doc:`predicate`
**default**       nil
================= =================

The predicate for selection of element/s.

actions
-------

================= =================
**required**      false
**type**          [] :doc:`action`
**default**       nil
================= =================

The actions to perform on the selected elements or page.
