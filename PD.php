<?php
/**
 * PD
 * ProxyDebug
 *
 * @author    kun <yangrokety@gmail.com>
 * @copyright 2014 kun
 * @license   http://www.php.net/license/3_01.txt  PHP License 3.01
 * @version   1.0
 * @link      https://github.com/yangxikun/tag-parse
 * @since     1.0
 */
/**
* PD
*
* @author    rokety <yangrokety@gmail.com>
* @license   http://www.php.net/license/3_01.txt  PHP License 3.01
* @version   1.0
* @link
* @since     1.0
*/
class PD
{
    protected static $debugItemCount = 0;
    protected static $debugGroupCount = 0;
    protected static $group = array();
    protected static $start;
    protected static $offset = 0;
    protected static $varNameString = null;

    /**
     * getVarName
     * get the variable name
     *
     * @access protected
     * @static
     *
     * @return string
     */
    protected static function getVarName()
    {
        if (self::$varNameString === null) {
            $trace = debug_backtrace();
            $line = file($trace[3]['file']);
            self::$varNameString = $line[$trace[3]['line']-1];
        }

        preg_match(
            '~\$([\w\d_]+)~',
            self::$varNameString,
            $matches,
            PREG_OFFSET_CAPTURE,
            self::$offset
        );

        if (!isset($matches[1])) {
            throw new Exception('Error Params, should use $variable as params', 1);
        }
        self::$offset = $matches[1][1];

        return $matches[1][0];
    }

    /**
     * func
     *
     * @param string $type debug type(info, warn, error)
     * @param mixed  $arg  debug variable
     *
     * @access protected
     * @static
     *
     * @return null
     */
    protected static function func($type, $arg)
    {
        if (self::$start) {
            self::$group[] = array(
                "category"=>$type,
                "type"=>gettype($arg),
                "name"=>self::getVarName(),
                "value"=>$arg
            );
        } else {
            self::$debugItemCount++;
            header(
                'Proxy_debug_item_'.self::$debugItemCount.': '
                .json_encode(
                    ["category"=>$type,
                    "type"=>gettype($arg),
                    "name"=>self::getVarName(),
                    "value"=>$arg]
                )
            );
            header('Proxy_debug_item_count: '.self::$debugItemCount);
        }
    }

    public static function __callStatic($name, $args)
    {
        $func = ['info'=>'I', 'warn'=>'W', 'error'=>'E'];
        if (isset($func[$name])) {
            self::$offset = 0;
            self::$varNameString = null;
            foreach ($args as $key => $arg) {
                self::func($func[$name], $arg);
            }
        } else {
            throw new Exception('Call to undefined method!', 1);
        }
    }

    /**
     * groupStart
     * start record a group
     *
     * @access public
     * @static
     *
     * @return null
     */
    public static function groupStart()
    {
        self::$start = true;
        self::$debugGroupCount++;
    }

    /**
     * groupEnd
     * stop record a group
     *
     * @access public
     * @static
     *
     * @return null
     */
    public static function groupEnd()
    {
        self::$start = false;
        header(
            'Proxy_debug_group_'
            .self::$debugGroupCount
            .': '.json_encode(self::$group)
        );
        header('Proxy_debug_group_count: '.self::$debugGroupCount);
        self::$group = array();
    }
}
